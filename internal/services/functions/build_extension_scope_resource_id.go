package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
)

// ExtensionResourceIdFunction builds resource IDs for extensions
type ExtensionResourceIdFunction struct{}

var _ function.Function = &ExtensionResourceIdFunction{}
var _ tffwdocs.FunctionWithRenderOption = &ExtensionResourceIdFunction{}

func (f *ExtensionResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "extension_resource_id"
}

func (f *ExtensionResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "base_resource_id",
				Description:         "The base resource ID of the Azure resource.",
				MarkdownDescription: "The base resource ID of the Azure resource.",
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
				Description:         "The list of resource names to construct the extension resource ID.",
				MarkdownDescription: "The list of resource names to construct the extension resource ID.",
			},
		},
		Return:              function.StringReturn{},
		Summary:             "Builds an extension resource ID.",
		Description:         "This function constructs an Azure extension resource ID given the base resource ID, resource type, and additional resource names.",
		MarkdownDescription: "This function constructs an Azure extension resource ID given the base resource ID, resource type, and additional resource names.",
	}
}

func (f *ExtensionResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var baseResourceID types.String
	var resourceType types.String
	var resourceNamesParam types.List

	if response.Error = request.Arguments.Get(ctx, &baseResourceID, &resourceType, &resourceNamesParam); response.Error != nil {
		return
	}

	resourceNames := utils.ParseResourceNames(resourceNamesParam)

	resourceID, err := parse.NewResourceIDWithNestedResourceNames(
		resourceNames,
		baseResourceID.ValueString(),
		resourceType.ValueString(),
	)
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	response.Error = response.Result.Set(ctx, types.StringValue(resourceID.AzureResourceId))
}

func (f *ExtensionResourceIdFunction) RenderOption() tffwdocs.FunctionRenderOption {
	return tffwdocs.FunctionRenderOption{
		Examples: []tffwdocs.Example{
			{
				HCL: `
locals {
  resource_id    = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  resource_type  = "Microsoft.Authorization/locks"
  resource_names = ["mylock"]
}

// it will output "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/providers/Microsoft.Authorization/locks/mylock"
output "extension_resource_id" {
  value = provider::azapi::extension_resource_id(local.resource_id, local.resource_type, local.resource_names)
}
`,
			},
		},
	}
}
