package services_test

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAzapiResourceReadSetsIdentityWhenArmReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	resourceID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest-rg"
	resourceType := "Microsoft.Resources/resourceGroups@2021-04-01"

	resourceClient, err := clients.NewResourceClient(staticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: staticTransport{
				statusCode: http.StatusNotFound,
				body:       `{}`,
			},
		},
	})
	if err != nil {
		t.Fatalf("building resource client: %+v", err)
	}

	azapiResource := &services.AzapiResource{
		ProviderData: &clients.Client{
			ResourceClient: resourceClient,
		},
	}

	var schemaResponse resource.SchemaResponse
	azapiResource.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)

	var identitySchemaResponse resource.IdentitySchemaResponse
	azapiResource.IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &identitySchemaResponse)

	model := services.NewDefaultAzapiResourceModel()
	model.ID = types.StringValue(resourceID)
	model.Type = types.StringValue(resourceType)
	model.Name = types.StringValue("acctest-rg")
	model.Location = types.StringValue("eastus")
	model.Body = types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}))

	state := tfsdk.State{Schema: schemaResponse.Schema}
	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("setting state: %+v", diags)
	}

	requestIdentity := &tfsdk.ResourceIdentity{
		Schema: identitySchemaResponse.IdentitySchema,
		Raw:    tftypes.NewValue(identitySchemaResponse.IdentitySchema.Type().TerraformType(ctx), nil),
	}
	responseIdentity := &tfsdk.ResourceIdentity{
		Schema: identitySchemaResponse.IdentitySchema,
		Raw:    tftypes.NewValue(identitySchemaResponse.IdentitySchema.Type().TerraformType(ctx), nil),
	}

	readRequest := resource.ReadRequest{
		State:    state,
		Identity: requestIdentity,
	}
	readResponse := resource.ReadResponse{
		State:    state,
		Identity: responseIdentity,
	}
	setEmptyPrivateState(&readRequest)
	setEmptyPrivateState(&readResponse)

	azapiResource.Read(ctx, readRequest, &readResponse)
	if readResponse.Diagnostics.HasError() {
		t.Fatalf("reading resource: %+v", readResponse.Diagnostics)
	}
	if readResponse.Identity.Raw.IsFullyNull() {
		t.Fatal("expected identity to be populated from the existing resource state")
	}

	var identity services.AzapiResourceIdentityModel
	diags = readResponse.Identity.Get(ctx, &identity)
	if diags.HasError() {
		t.Fatalf("reading identity: %+v", diags)
	}

	if identity.ID.ValueString() != resourceID {
		t.Fatalf("expected identity ID %q, got %q", resourceID, identity.ID.ValueString())
	}
}

func TestAzapiResourceMoveStateSetsTargetIdentityForAzurermResources(t *testing.T) {
	ctx := context.Background()
	sourceID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest-rg"

	azapiResource := &services.AzapiResource{}
	stateMovers := azapiResource.MoveState(ctx)
	if len(stateMovers) == 0 {
		t.Fatal("expected at least one state mover")
	}

	stateMover := stateMovers[0]

	sourceState := tfsdk.State{Schema: *stateMover.SourceSchema}
	diags := sourceState.Set(ctx, azurermResourceGroupMoveState{
		ID:                    types.StringValue(sourceID),
		ResourceManagerID:     types.StringNull(),
		ResourceID:            types.StringNull(),
		ResourceVersionlessID: types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("setting source state: %+v", diags)
	}

	var schemaResponse resource.SchemaResponse
	azapiResource.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("reading azapi schema: %+v", schemaResponse.Diagnostics)
	}

	var identitySchemaResponse resource.IdentitySchemaResponse
	azapiResource.IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &identitySchemaResponse)
	if identitySchemaResponse.Diagnostics.HasError() {
		t.Fatalf("reading azapi identity schema: %+v", identitySchemaResponse.Diagnostics)
	}

	moveRequest := resource.MoveStateRequest{
		SourceTypeName: "azurerm_resource_group",
		SourceState:    &sourceState,
	}

	moveResponse := resource.MoveStateResponse{
		TargetState: tfsdk.State{
			Schema: schemaResponse.Schema,
		},
		TargetIdentity: &tfsdk.ResourceIdentity{
			Schema: identitySchemaResponse.IdentitySchema,
			Raw:    tftypes.NewValue(identitySchemaResponse.IdentitySchema.Type().TerraformType(ctx), nil),
		},
	}
	setEmptyMoveStateTargetPrivate(&moveResponse)

	stateMover.StateMover(ctx, moveRequest, &moveResponse)
	if moveResponse.Diagnostics.HasError() {
		t.Fatalf("moving state: %+v", moveResponse.Diagnostics)
	}

	var movedState services.AzapiResourceModel
	diags = moveResponse.TargetState.Get(ctx, &movedState)
	if diags.HasError() {
		t.Fatalf("reading moved state: %+v", diags)
	}

	var movedIdentity services.AzapiResourceIdentityModel
	diags = moveResponse.TargetIdentity.Get(ctx, &movedIdentity)
	if diags.HasError() {
		t.Fatalf("reading moved identity: %+v", diags)
	}

	if movedIdentity.ID.ValueString() != movedState.ID.ValueString() {
		t.Fatalf("expected moved identity ID %q, got %q", movedState.ID.ValueString(), movedIdentity.ID.ValueString())
	}
	if movedIdentity.Type.ValueString() != movedState.Type.ValueString() {
		t.Fatalf("expected moved identity type %q, got %q", movedState.Type.ValueString(), movedIdentity.Type.ValueString())
	}
}

type staticTokenCredential struct{}

func (staticTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

type staticTransport struct {
	statusCode int
	body       string
}

func (t staticTransport) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: t.statusCode,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(t.body)),
		Request: req,
	}, nil
}

func setEmptyPrivateState(target any) {
	field := reflect.ValueOf(target).Elem().FieldByName("Private")
	field.Set(reflect.New(field.Type().Elem()))
}

func setEmptyMoveStateTargetPrivate(target any) {
	field := reflect.ValueOf(target).Elem().FieldByName("TargetPrivate")
	field.Set(reflect.New(field.Type().Elem()))
}

type azurermResourceGroupMoveState struct {
	ID                    types.String `tfsdk:"id"`
	ResourceManagerID     types.String `tfsdk:"resource_manager_id"`
	ResourceID            types.String `tfsdk:"resource_id"`
	ResourceVersionlessID types.String `tfsdk:"resource_versionless_id"`
}
