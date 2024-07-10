package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuildTenantScopeResourceIdFunction struct{}

func (t *BuildTenantScopeResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "build_tenant_scope_resource_id"
}

func (t *BuildTenantScopeResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
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

func (t *BuildTenantScopeResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var resourceType types.String
	var name types.String

	if response.Error = request.Arguments.Get(ctx, &resourceType, &name); response.Error != nil {
		return
	}

	parentID := "/"

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

var _ function.Function = &BuildTenantScopeResourceIdFunction{}
