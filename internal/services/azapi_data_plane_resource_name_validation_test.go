package services

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateDataPlaneResourceAddress(t *testing.T) {
	t.Run("agent requires name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/agents@v1"),
			Name: types.StringNull(),
		}

		err := validateDataPlaneResourceAddress(config)
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "must be set") {
			t.Fatalf("expected must-be-set error, got: %v", err)
		}
	})

	t.Run("agent accepts explicit name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/agents@v1"),
			Name: types.StringValue("terraform-agent"),
		}

		if err := validateDataPlaneResourceAddress(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("non-assistant requires name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.KeyVault/vaults/secrets@7.4"),
			Name: types.StringNull(),
		}

		err := validateDataPlaneResourceAddress(config)
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "must be set") {
			t.Fatalf("expected must-be-set error, got: %v", err)
		}
	})

	t.Run("non-assistant accepts name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.KeyVault/vaults/secrets@7.4"),
			Name: types.StringValue("secret-name"),
		}

		if err := validateDataPlaneResourceAddress(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("table entity requires identifiers instead of name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"),
			Name: types.StringNull(),
			Identifiers: types.MapValueMust(types.StringType, map[string]attr.Value{
				"partitionKey": types.StringValue("pk"),
				"rowKey":       types.StringValue("rk"),
			}),
		}

		if err := validateDataPlaneResourceAddress(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("table entity rejects missing identifiers", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type:        types.StringValue("Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"),
			Name:        types.StringNull(),
			Identifiers: types.MapNull(types.StringType),
		}

		err := validateDataPlaneResourceAddress(config)
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "identifiers") {
			t.Fatalf("expected identifiers error, got: %v", err)
		}
	})

	t.Run("table entity skips validation when identifier values are unknown (pre-for_each expansion)", func(t *testing.T) {
		// Terraform calls ValidateConfig once on the raw block config before expanding
		// for_each.  At that point each.value.pk / each.value.rk are unknown string values
		// inside an otherwise-known map.  Validation must be deferred, not rejected.
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"),
			Name: types.StringNull(),
			Identifiers: types.MapValueMust(types.StringType, map[string]attr.Value{
				"partitionKey": types.StringUnknown(),
				"rowKey":       types.StringUnknown(),
			}),
		}

		if err := validateDataPlaneResourceAddress(config); err != nil {
			t.Fatalf("expected nil (deferred) error for unknown identifier values, got: %v", err)
		}
	})

	t.Run("table entity skips validation when any single identifier value is unknown", func(t *testing.T) {
		// Mixed case: one key has a concrete value, the other is still unknown.
		// AsMapOfString returns an empty map when ElementsAs encounters any unknown
		// element (allowUnknowns=false), so ALL keys appear missing — not just the
		// unknown one.  We must defer when ANY element is unknown to avoid a false error.
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"),
			Name: types.StringNull(),
			Identifiers: types.MapValueMust(types.StringType, map[string]attr.Value{
				"partitionKey": types.StringValue("connectivity-hub"),
				"rowKey":       types.StringUnknown(),
			}),
		}

		if err := validateDataPlaneResourceAddress(config); err != nil {
			t.Fatalf("expected nil (deferred) error for mixed known/unknown identifier values, got: %v", err)
		}
	})
}
