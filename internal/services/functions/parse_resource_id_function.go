package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
	"github.com/magodo/terraform-plugin-framework-docs/fwdtypes"
)

type ParseResourceIdFunction struct {
}

var _ function.Function = &ParseResourceIdFunction{}
var _ tffwdocs.FunctionWithRenderOption = &ParseResourceIdFunction{}

func ParseResourceIdResultAttrTypes(withDoc bool) map[string]attr.Type {
	if withDoc {
		return map[string]attr.Type{
			"id":                  fwdtypes.NewStringType("The resource id of this resource."),
			"type":                fwdtypes.NewStringType("The azure resource type."),
			"name":                fwdtypes.NewStringType("The resource name."),
			"parent_id":           fwdtypes.NewStringType("The resource id of the parent resource."),
			"resource_group_name": fwdtypes.NewStringType("The name of the resource group this resource resides in."),
			"resource_group_id":   fwdtypes.NewStringType("The id of the resource group this resource resides in."),
			"subscription_id":     fwdtypes.NewStringType("The id of the subscription this resource resides in."),
			"provider_namespace":  fwdtypes.NewStringType("The namespace of the resource provider."),
			"parts":               fwdtypes.NewMapType("A map of the parts of the resource id.", types.StringType),
		}
	} else {
		return map[string]attr.Type{
			"id":                  types.StringType,
			"type":                types.StringType,
			"name":                types.StringType,
			"parent_id":           types.StringType,
			"resource_group_name": types.StringType,
			"resource_group_id":   types.StringType,
			"subscription_id":     types.StringType,
			"provider_namespace":  types.StringType,
			"parts":               types.MapType{ElemType: types.StringType},
		}
	}
}

func (p *ParseResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "parse_resource_id"
}

func (p *ParseResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "resource_type",
				Description:         "The resource type of the Azure resource.",
				MarkdownDescription: "The resource type of the Azure resource.",
			},
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "resource_id",
				Description:         "The resource ID of the Azure resource to parse.",
				MarkdownDescription: "The resource ID of the Azure resource to parse.",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: ParseResourceIdResultAttrTypes(true),
		},
		Summary:             "Parses an Azure resource ID into its components.",
		Description:         "This function takes an Azure resource ID and a resource type and parses the ID into its individual components such as subscription ID, resource group name, provider namespace, and other parts.",
		MarkdownDescription: "This function takes an Azure resource ID and a resource type and parses the ID into its individual components such as subscription ID, resource group name, provider namespace, and other parts.",
		DeprecationMessage:  "",
	}
}

func (p *ParseResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var resourceTypeParam string
	var resourceId types.String

	if response.Error = request.Arguments.Get(ctx, &resourceTypeParam, &resourceId); response.Error != nil {
		return
	}

	resourceType := utils.TryAppendDefaultApiVersion(resourceTypeParam)

	id, err := parse.ResourceIDWithResourceType(resourceId.ValueString(), resourceType)
	if err != nil {
		response.Error = function.NewFuncError(fmt.Errorf("failed to parse resource ID(resourceType: %s, resourceId: %s): %w", resourceType, resourceId.ValueString(), err).Error())
		return
	}

	armId, err := arm.ParseResourceID(id.AzureResourceId)
	if id.AzureResourceId == "/" {
		armId, err = &arm.ResourceID{
			ResourceType: arm.TenantResourceType,
		}, nil
	}
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	path := id.AzureResourceId
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	components := strings.Split(path, "/")
	parts := make(map[string]attr.Value)
	for i := 0; i < len(components)-1; i += 2 {
		parts[components[i]] = basetypes.NewStringValue(components[i+1])
	}

	resourceGroupId := ""
	if armId.ResourceGroupName != "" {
		resourceGroupId = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", armId.SubscriptionID, armId.ResourceGroupName)
	}

	result := map[string]attr.Value{
		"id":                  types.StringValue(id.ID()),
		"type":                types.StringValue(id.AzureResourceType),
		"name":                types.StringValue(id.Name),
		"parent_id":           types.StringValue(id.ParentId),
		"resource_group_name": types.StringValue(armId.ResourceGroupName),
		"resource_group_id":   types.StringValue(resourceGroupId),
		"subscription_id":     types.StringValue(armId.SubscriptionID),
		"provider_namespace":  types.StringValue(armId.ResourceType.Namespace),
		"parts":               types.MapValueMust(types.StringType, parts),
	}

	response.Error = response.Result.Set(ctx, types.ObjectValueMust(ParseResourceIdResultAttrTypes(false), result))
}

func (p *ParseResourceIdFunction) RenderOption() tffwdocs.FunctionRenderOption {
	return tffwdocs.FunctionRenderOption{
		Examples: []tffwdocs.Example{
			{
				HCL: `
locals {
  resource_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
  resource_type = "Microsoft.Network/virtualNetworks"
}

// it will output below object
# {
#   "id" = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
#   "name" = "myVNet"
#   "parent_id" = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
#   "parts" = tomap({
#     "providers" = "Microsoft.Network"
#     "resourceGroups" = "myResourceGroup"
#     "subscriptions" = "00000000-0000-0000-0000-000000000000"
#     "virtualNetworks" = "myVNet"
#   })
#   "provider_namespace" = "Microsoft.Network"
#   "resource_group_id" = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
#   "resource_group_name" = "myResourceGroup"
#   "subscription_id" = "00000000-0000-0000-0000-000000000000"
#   "type" = "Microsoft.Network/virtualNetworks"
# }
output "parsed_resource_id" {
  value = provider::azapi::parse_resource_id(local.resource_type, local.resource_id)
}
`,
			},
		},
	}
}
