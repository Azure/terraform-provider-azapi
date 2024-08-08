package functions

import (
	"context"
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceGroupResourceIdFunction builds resource IDs for resource group scope
type ResourceGroupResourceIdFunction struct{}

func (f *ResourceGroupResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "resource_group_resource_id"
}

func (f *ResourceGroupResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "subscription_id",
			},
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "resource_group_name",
			},
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "resource_type",
			},
			function.ListParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "resource_names",
				ElementType:        types.StringType,
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: BuildResourceIdResultAttrTypes,
		},
		Summary:             "Builds a resource group scope resource ID.",
		Description:         "This function constructs an Azure resource group scope resource ID given the subscription ID, resource group name, resource type, and resource names.",
		MarkdownDescription: "This function constructs an Azure resource group scope resource ID given the subscription ID, resource group name, resource type, and resource names.",
		DeprecationMessage:  "",
	}
}

func (f *ResourceGroupResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var subscriptionID types.String
	var resourceGroupName types.String
	var resourceType types.String
	var resourceNamesParam types.List

	if response.Error = request.Arguments.Get(ctx, &subscriptionID, &resourceGroupName, &resourceType, &resourceNamesParam); response.Error != nil {
		return
	}

	resourceNames := utils.ParseResourceNames(resourceNamesParam)

	resourceID, err := parse.NewResourceIDWithNestedResourceNames(
		resourceNames,
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID.ValueString(), resourceGroupName.ValueString()),
		resourceType.ValueString(),
	)
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	result := map[string]attr.Value{
		"resource_id": types.StringValue(resourceID.AzureResourceId),
	}

	response.Error = response.Result.Set(ctx, types.ObjectValueMust(BuildResourceIdResultAttrTypes, result))
}

var _ function.Function = &ResourceGroupResourceIdFunction{}
