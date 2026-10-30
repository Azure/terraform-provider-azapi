package services_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func dynamicFromJSON(t *testing.T, input string) types.Dynamic {
	t.Helper()
	out, err := dynamic.FromJSONImplied([]byte(input))
	require.NoError(t, err)
	return out
}

func newDataPlaneResourceWithTransport(t *testing.T, transport policy.Transporter) *services.DataPlaneResource {
	t.Helper()
	dataPlaneClient, err := clients.NewDataPlaneClient(staticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	require.NoError(t, err)

	return &services.DataPlaneResource{
		ProviderData: &clients.Client{
			Features:        features.Default(),
			DataPlaneClient: dataPlaneClient,
		},
	}
}

func dataPlaneResourceSchema(t *testing.T, r *services.DataPlaneResource) schema.Schema {
	t.Helper()
	var schemaResponse resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResponse)
	require.False(t, schemaResponse.Diagnostics.HasError(), "%+v", schemaResponse.Diagnostics)
	return schemaResponse.Schema
}

func newDataPlaneResourceModel(resourceType, parentID, name string, body types.Dynamic) services.DataPlaneResourceModel {
	nameValue := types.StringNull()
	if name != "" {
		nameValue = types.StringValue(name)
	}
	id, _ := parse.NewDataPlaneResourceId(name, parentID, resourceType)
	return services.DataPlaneResourceModel{
		ID:                            types.StringValue(id.ID()),
		Name:                          nameValue,
		ParentID:                      types.StringValue(parentID),
		Type:                          types.StringValue(resourceType),
		Body:                          body,
		SensitiveBody:                 types.DynamicNull(),
		SensitiveBodyVersion:          types.MapNull(types.StringType),
		IgnoreCasing:                  types.BoolValue(false),
		IgnoreMissingProperty:         types.BoolValue(true),
		ReplaceTriggersExternalValues: types.DynamicNull(),
		ReplaceTriggersRefs:           types.ListNull(types.StringType),
		ResponseExportValues:          types.DynamicNull(),
		Retry:                         retry.NewRetryValueNull(),
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
}
