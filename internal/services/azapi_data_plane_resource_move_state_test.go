package services_test

import (
	"context"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type azurermKeyVaultSecretMoveState struct {
	ID            types.String `tfsdk:"id"`
	VersionlessID types.String `tfsdk:"versionless_id"`
}

func moveToDataPlaneResource(t *testing.T, sourceProviderAddress, sourceTypeName string, source azurermKeyVaultSecretMoveState) resource.MoveStateResponse {
	t.Helper()
	ctx := context.Background()

	r := &services.DataPlaneResource{}
	stateMovers := r.MoveState(ctx)
	if len(stateMovers) != 1 {
		t.Fatalf("expected one state mover, got %d", len(stateMovers))
	}

	sourceState := tfsdk.State{Schema: *stateMovers[0].SourceSchema}
	if diags := sourceState.Set(ctx, source); diags.HasError() {
		t.Fatalf("setting source state: %+v", diags)
	}

	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("reading schema: %+v", schemaResponse.Diagnostics)
	}

	// The framework hands each mover a null target state, and treats a mover that leaves it null as
	// not having handled the move.
	response := resource.MoveStateResponse{
		TargetState: tfsdk.State{
			Schema: schemaResponse.Schema,
			Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
		},
	}
	stateMovers[0].StateMover(ctx, resource.MoveStateRequest{
		SourceProviderAddress: sourceProviderAddress,
		SourceTypeName:        sourceTypeName,
		SourceState:           &sourceState,
	}, &response)

	return response
}

func TestDataPlaneResourceMoveState_keyVaultSecret(t *testing.T) {
	testCases := map[string]struct {
		source         azurermKeyVaultSecretMoveState
		expectedParent string
	}{
		"versionless id": {
			source: azurermKeyVaultSecretMoveState{
				ID:            types.StringValue("https://myvault.vault.azure.net/secrets/mysecret/0123456789abcdef0123456789abcdef"),
				VersionlessID: types.StringValue("https://myvault.vault.azure.net/secrets/mysecret"),
			},
			expectedParent: "myvault.vault.azure.net",
		},
		"versioned id only": {
			source: azurermKeyVaultSecretMoveState{
				ID:            types.StringValue("https://myvault.vault.azure.net/secrets/mysecret/0123456789abcdef0123456789abcdef"),
				VersionlessID: types.StringNull(),
			},
			expectedParent: "myvault.vault.azure.net",
		},
		"sovereign cloud": {
			source: azurermKeyVaultSecretMoveState{
				ID:            types.StringValue("https://myvault.vault.usgovcloudapi.net/secrets/mysecret/0123456789abcdef0123456789abcdef"),
				VersionlessID: types.StringValue("https://myvault.vault.usgovcloudapi.net/secrets/mysecret"),
			},
			expectedParent: "myvault.vault.usgovcloudapi.net",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			response := moveToDataPlaneResource(t, "registry.terraform.io/hashicorp/azurerm", "azurerm_key_vault_secret", tc.source)
			if response.Diagnostics.HasError() {
				t.Fatalf("moving state: %+v", response.Diagnostics)
			}

			expected := map[string]string{
				"id":        tc.expectedParent + "/secrets/mysecret",
				"name":      "mysecret",
				"parent_id": tc.expectedParent,
				"type":      "Microsoft.KeyVault/vaults/secrets@7.5",
			}
			for attribute, want := range expected {
				var got types.String
				if diags := response.TargetState.GetAttribute(ctx, path.Root(attribute), &got); diags.HasError() {
					t.Fatalf("reading %s: %+v", attribute, diags)
				}
				if got.ValueString() != want {
					t.Fatalf("expected %s %q, got %q", attribute, want, got.ValueString())
				}
			}

			// The moved state must match an imported one, which never holds the secret value.
			for _, attribute := range []string{"body", "sensitive_body"} {
				var got types.Dynamic
				if diags := response.TargetState.GetAttribute(ctx, path.Root(attribute), &got); diags.HasError() {
					t.Fatalf("reading %s: %+v", attribute, diags)
				}
				if !got.IsNull() {
					t.Fatalf("expected %s to be null, got %s", attribute, got.String())
				}
			}
		})
	}
}

func TestDataPlaneResourceMoveState_otherSourcesAreNotHandled(t *testing.T) {
	source := azurermKeyVaultSecretMoveState{
		ID:            types.StringValue("https://myvault.vault.azure.net/keys/mykey/0123456789abcdef0123456789abcdef"),
		VersionlessID: types.StringValue("https://myvault.vault.azure.net/keys/mykey"),
	}

	testCases := map[string]struct {
		sourceProviderAddress string
		sourceTypeName        string
	}{
		"unsupported azurerm type": {
			sourceProviderAddress: "registry.terraform.io/hashicorp/azurerm",
			sourceTypeName:        "azurerm_key_vault_key",
		},
		"another provider": {
			sourceProviderAddress: "registry.terraform.io/example/azurerm",
			sourceTypeName:        "azurerm_key_vault_secret",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			response := moveToDataPlaneResource(t, tc.sourceProviderAddress, tc.sourceTypeName, source)
			if response.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics, got %+v", response.Diagnostics)
			}
			if !response.TargetState.Raw.IsNull() {
				t.Fatalf("expected the target state to be left null so the move is not claimed, got %s", response.TargetState.Raw.String())
			}
		})
	}
}

func TestDataPlaneResourceMoveState_invalidSourceID(t *testing.T) {
	response := moveToDataPlaneResource(t, "registry.terraform.io/hashicorp/azurerm", "azurerm_key_vault_secret", azurermKeyVaultSecretMoveState{
		ID:            types.StringValue("not-a-key-vault-url"),
		VersionlessID: types.StringNull(),
	})
	if !response.Diagnostics.HasError() {
		t.Fatal("expected an error for a source ID that is not a Key Vault object URL")
	}
}
