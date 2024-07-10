package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ParseResourceIdFunction struct {
}

var ParseResourceIdResultAttrTypes = map[string]attr.Type{
	"id":                  types.StringType,
	"type":                types.StringType,
	"name":                types.StringType,
	"parent_id":           types.StringType,
	"resource_group_name": types.StringType,
	"subscription_id":     types.StringType,
	"provider_namespace":  types.StringType,
	"parts": types.MapType{
		ElemType: types.StringType,
	},
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
				CustomType:          nil,
				Description:         "",
				MarkdownDescription: "",
				Name:                "type",
			},
			function.StringParameter{
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				CustomType:          nil,
				Description:         "",
				MarkdownDescription: "",
				Name:                "resource_id",
			},
		},
		VariadicParameter: nil,
		Return: function.ObjectReturn{
			AttributeTypes: ParseResourceIdResultAttrTypes,
		},
		Summary:             "",
		Description:         "",
		MarkdownDescription: "",
		DeprecationMessage:  "",
	}
}

func (p *ParseResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var resourceType string
	var resourceId types.String

	if response.Error = request.Arguments.Get(ctx, &resourceType, &resourceId); response.Error != nil {
		return
	}

	if resourceId.IsNull() {
		response.Error = response.Result.Set(ctx, types.ObjectNull(ParseResourceIdResultAttrTypes))
		return
	}

	if resourceId.IsUnknown() {
		response.Error = response.Result.Set(ctx, types.ObjectUnknown(ParseResourceIdResultAttrTypes))
		return
	}

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

	result := map[string]attr.Value{
		"id":                  types.StringValue(id.ID()),
		"type":                types.StringValue(id.AzureResourceType),
		"name":                types.StringValue(id.Name),
		"parent_id":           types.StringValue(id.ParentId),
		"resource_group_name": types.StringValue(armId.ResourceGroupName),
		"subscription_id":     types.StringValue(armId.SubscriptionID),
		"provider_namespace":  types.StringValue(armId.ResourceType.Namespace),
		"parts":               types.MapValueMust(types.StringType, parts),
	}

	response.Error = response.Result.Set(ctx, types.ObjectValueMust(ParseResourceIdResultAttrTypes, result))
}

var _ function.Function = &ParseResourceIdFunction{}
