package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TenantResourceIdFunction builds resource IDs for tenant scope
type TenantResourceIdFunction struct{}

func (f *TenantResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "tenant_resource_id"
}

func (f *TenantResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
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
		Summary:             "Builds a tenant scope resource ID.",
		Description:         "This function constructs an Azure tenant scope resource ID given the resource type and resource names.",
		MarkdownDescription: "This function constructs an Azure tenant scope resource ID given the resource type and resource names.",
	}
}

func (f *TenantResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var resourceType types.String
	var resourceNamesParam types.List

	if response.Error = request.Arguments.Get(ctx, &resourceType, &resourceNamesParam); response.Error != nil {
		return
	}

	resourceNames := utils.ParseResourceNames(resourceNamesParam)

	resourceID, err := parse.NewResourceIDWithNestedResourceNames(resourceNames, "/", resourceType.ValueString())
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	response.Error = response.Result.Set(ctx, types.StringValue(resourceID.AzureResourceId))
}

var _ function.Function = &TenantResourceIdFunction{}
