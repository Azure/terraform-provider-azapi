package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuildManagementGroupScopeResourceIdFunction struct{}

func (m *BuildManagementGroupScopeResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "build_management_group_scope_resource_id"
}

func (m *BuildManagementGroupScopeResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "management_group_id",
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

func (m *BuildManagementGroupScopeResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var mgID types.String
	var resourceType types.String
	var name types.String

	if response.Error = request.Arguments.Get(ctx, &mgID, &resourceType, &name); response.Error != nil {
		return
	}

	if mgID.ValueString() == "" {
		response.Error = function.NewFuncError("management_group_id cannot be empty")
		return
	}

	parentID := "/providers/Microsoft.Management/managementGroups/" + mgID.ValueString()

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

var _ function.Function = &BuildManagementGroupScopeResourceIdFunction{}
