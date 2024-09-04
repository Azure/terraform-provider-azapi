package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
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
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "management_group_name",
				Description:         "The name of the management group.",
				MarkdownDescription: "The name of the management group.",
			},
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "resource_type",
				Description:         "The resource type of the Azure resource.",
				MarkdownDescription: "The resource type of the Azure resource.",
			},
			function.ListParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "resource_names",
				ElementType:         types.StringType,
				Description:         "The list of resource names to construct the resource ID.",
				MarkdownDescription: "The list of resource names to construct the resource ID.",
			},
		},
		Return:              function.StringReturn{},
		Summary:             "Builds a management group scope resource ID.",
		Description:         "This function constructs an Azure management group scope resource ID given the management group name, resource type, and resource names.",
		MarkdownDescription: "This function constructs an Azure management group scope resource ID given the management group name, resource type, and resource names.",
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

	response.Error = response.Result.Set(ctx, types.StringValue(resourceID.AzureResourceId))
}

var _ function.Function = &ManagementGroupResourceIdFunction{}
