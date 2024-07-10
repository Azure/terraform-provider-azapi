package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuildResourceIdFunction struct{}

var BuildResourceIdResultAttrTypes = map[string]attr.Type{
	"resource_id": types.StringType,
}

func (b *BuildResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "build_resource_id"
}

func (b *BuildResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:     false,
				AllowUnknownValues: false,
				Name:               "parent_id",
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

func (b *BuildResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var resourceType types.String
	var parentId types.String
	var name types.String

	if response.Error = request.Arguments.Get(ctx, &parentId, &resourceType, &name); response.Error != nil {
		return
	}

	resourceID, err := parse.NewResourceID(name.ValueString(), parentId.ValueString(), resourceType.ValueString())
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	result := map[string]attr.Value{
		"resource_id": types.StringValue(resourceID.ID()),
	}

	response.Error = response.Result.Set(ctx, types.ObjectValueMust(BuildResourceIdResultAttrTypes, result))
}

var _ function.Function = &BuildResourceIdFunction{}
