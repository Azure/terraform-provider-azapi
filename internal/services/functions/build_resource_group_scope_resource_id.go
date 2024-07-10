package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuildResourceGroupScopeResourceIdFunction struct{}

func (r *BuildResourceGroupScopeResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "build_resource_group_scope_resource_id"
}

func (r *BuildResourceGroupScopeResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
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
				Name:               "type",
			},
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "name",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: BuildResourceIdResultAttrTypes,
		},
	}
}

func (r *BuildResourceGroupScopeResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var subID types.String
	var rgName types.String
	var resourceType types.String
	var name types.String

	if response.Error = request.Arguments.Get(ctx, &subID, &rgName, &resourceType, &name); response.Error != nil {
		return
	}

	if subID.ValueString() == "" {
		response.Error = function.NewFuncError("subscription_id cannot be empty")
		return
	}

	if rgName.ValueString() == "" {
		response.Error = function.NewFuncError("resource_group_name cannot be empty")
		return
	}

	parentID := "/subscriptions/" + subID.ValueString() + "/resourceGroups/" + rgName.ValueString()

	resourceID, err := parse.NewResourceID(name.ValueString(), parentID, resourceType.ValueString())
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	result := map[string]attr.Value{
		"resource_id": types.StringValue(resourceID.ID()),
	}

	response.Error = response.Result.Set(ctx, types.ObjectValueMust(BuildResourceIdResultAttrTypes, result))
}

var _ function.Function = &BuildResourceGroupScopeResourceIdFunction{}
