package services_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDataPlaneResourceCreateRetriesAuthorizationFailureDuringExistenceCheck(t *testing.T) {
	transport := &sequenceTransport{responses: []sequenceResponse{
		{
			method:     http.MethodGet,
			statusCode: http.StatusForbidden,
			body:       `{"error":{"code":"Forbidden","message":"role assignment has not propagated"}}`,
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
			body:       `{"error":{"code":"ResourceNotFound","message":"resource does not exist"}}`,
		},
		{
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       dataPlaneResourceResponse,
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       dataPlaneResourceResponse,
		},
	}}

	diagnostics := createDataPlaneResourceWithRetry(t, transport, []string{"Forbidden"})
	if diagnostics.HasError() {
		t.Fatalf("creating data plane resource: %+v", diagnostics)
	}

	assertRequestMethods(t, transport, []string{http.MethodGet, http.MethodGet, http.MethodPut, http.MethodGet})
}

func TestDataPlaneResourceCreateDoesNotRetryNotFoundDuringExistenceCheck(t *testing.T) {
	transport := &sequenceTransport{responses: []sequenceResponse{
		{
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
			body:       `{"error":{"code":"FooResourceNotFound","message":"resource does not exist"}}`,
		},
		{
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       dataPlaneResourceResponse,
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       dataPlaneResourceResponse,
		},
	}}

	diagnostics := createDataPlaneResourceWithRetry(t, transport, []string{"FooResourceNotFound"})
	if diagnostics.HasError() {
		t.Fatalf("creating data plane resource: %+v", diagnostics)
	}

	assertRequestMethods(t, transport, []string{http.MethodGet, http.MethodPut, http.MethodGet})
}

const dataPlaneResourceResponse = `{
	"key":"test",
	"value":"value"
}`

func createDataPlaneResourceWithRetry(t *testing.T, transport policy.Transporter, regexes []string) diag.Diagnostics {
	t.Helper()

	ctx := context.Background()
	dataPlaneClient, err := clients.NewDataPlaneClient(staticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("building data plane client: %+v", err)
	}

	dataPlaneResource := &services.DataPlaneResource{
		ProviderData: &clients.Client{
			Features:        features.Default(),
			DataPlaneClient: dataPlaneClient,
		},
	}

	var schemaResponse resource.SchemaResponse
	dataPlaneResource.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("building data plane resource schema: %+v", schemaResponse.Diagnostics)
	}

	model := services.DataPlaneResourceModel{
		ID:                            types.StringNull(),
		Name:                          types.StringValue("test"),
		ParentID:                      types.StringValue("example.azconfig.io"),
		Type:                          types.StringValue("Microsoft.AppConfiguration/configurationStores/keyValues@1.0"),
		Body:                          types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{})),
		SensitiveBody:                 types.DynamicNull(),
		SensitiveBodyVersion:          types.MapNull(types.StringType),
		IgnoreCasing:                  types.BoolValue(false),
		IgnoreMissingProperty:         types.BoolValue(true),
		ReplaceTriggersExternalValues: types.DynamicNull(),
		ReplaceTriggersRefs:           types.ListNull(types.StringType),
		ResponseExportValues:          types.DynamicNull(),
		Retry:                         testRetryValue(regexes),
		Locks:                         types.ListNull(types.StringType),
		Output:                        types.DynamicNull(),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"read":   types.StringType,
				"delete": types.StringType,
			}),
		},
		CreateHeaders:         types.MapNull(types.StringType),
		CreateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		UpdateHeaders:         types.MapNull(types.StringType),
		UpdateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		DeleteHeaders:         types.MapNull(types.StringType),
		DeleteQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		ReadHeaders:           types.MapNull(types.StringType),
		ReadQueryParameters:   types.MapNull(types.ListType{ElemType: types.StringType}),
	}

	plan := tfsdk.Plan{Schema: schemaResponse.Schema}
	if diagnostics := plan.Set(ctx, &model); diagnostics.HasError() {
		t.Fatalf("setting data plane resource plan: %+v", diagnostics)
	}
	config := tfsdk.Config{
		Schema: schemaResponse.Schema,
		Raw:    plan.Raw,
	}

	state := tfsdk.State{
		Schema: schemaResponse.Schema,
		Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
	}

	var diagnostics diag.Diagnostics
	dataPlaneResource.CreateUpdate(ctx, config, plan, &state, &diagnostics, newTestPrivateData())
	return diagnostics
}
