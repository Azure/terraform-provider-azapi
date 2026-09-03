package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
)

// SubscriptionResourceIdFunction builds resource IDs for subscription scope
type SubscriptionResourceIdFunction struct{}

var _ function.Function = &SubscriptionResourceIdFunction{}
var _ tffwdocs.FunctionWithRenderOption = &SubscriptionResourceIdFunction{}

func (f *SubscriptionResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "subscription_resource_id"
}

func (f *SubscriptionResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "subscription_id",
				Description:         "The subscription ID of the Azure resource.",
				MarkdownDescription: "The subscription ID of the Azure resource.",
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
		Summary:             "Builds a subscription scope resource ID.",
		Description:         "This function constructs an Azure subscription scope resource ID given the subscription ID, resource type, and resource names.",
		MarkdownDescription: "This function constructs an Azure subscription scope resource ID given the subscription ID, resource type, and resource names.",
	}
}

func (f *SubscriptionResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var subscriptionID types.String
	var resourceType types.String
	var resourceNamesParam types.List

	if response.Error = request.Arguments.Get(ctx, &subscriptionID, &resourceType, &resourceNamesParam); response.Error != nil {
		return
	}

	resourceNames := utils.ParseResourceNames(resourceNamesParam)

	resourceID, err := parse.NewResourceIDWithNestedResourceNames(resourceNames, "/subscriptions/"+subscriptionID.ValueString(), resourceType.ValueString())
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	response.Error = response.Result.Set(ctx, types.StringValue(resourceID.AzureResourceId))
}

func (f *SubscriptionResourceIdFunction) RenderOption() tffwdocs.FunctionRenderOption {
	return tffwdocs.FunctionRenderOption{
		Examples: []tffwdocs.Example{
			{
				HCL: `
locals {
  subscription_id = "00000000-0000-0000-0000-000000000000"
  resource_type   = "Microsoft.Resources/resourceGroups"
  resource_names  = ["rg1"]
}

// it will output "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
output "subscription_resource_id" {
  value = provider::azapi::subscription_resource_id(local.subscription_id, local.resource_type, local.resource_names)
}
`,
			},
		},
	}
}
