package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ManagementGroupResourceIdFunction builds resource IDs for management group scope
type ManagementGroupResourceIdFunction struct{}

func (f *ManagementGroupResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "management_group_resource_id"
}

func (f *ManagementGroupResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "management_group_name",
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
		Summary:             "Builds a management group scope resource ID.",
		Description:         "This function constructs an Azure management group scope resource ID given the management group name, resource type, and resource names.",
		MarkdownDescription: "This function constructs an Azure management group scope resource ID given the management group name, resource type, and resource names.",
		DeprecationMessage:  "",
	}
}

func (f *ManagementGroupResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var managementGroupName types.String
	var resourceType types.String
	var resourceNamesParam types.List

	if response.Error = request.Arguments.Get(ctx, &managementGroupName, &resourceType, &resourceNamesParam); response.Error != nil {
		return
	}

	resourceNames := utils.ParseResourceNames(resourceNamesParam)

	resourceID, err := parse.NewResourceIDWithNestedResourceNames(resourceNames, "/providers/Microsoft.Management/managementGroups/"+managementGroupName.ValueString(), resourceType.ValueString())
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	result := map[string]attr.Value{
		"resource_id": types.StringValue(resourceID.AzureResourceId),
	}

	response.Error = response.Result.Set(ctx, types.ObjectValueMust(BuildResourceIdResultAttrTypes, result))
}

var _ function.Function = &ManagementGroupResourceIdFunction{}
