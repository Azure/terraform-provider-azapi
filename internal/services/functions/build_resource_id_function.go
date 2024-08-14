package functions

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuildResourceIdFunction struct{}

func (b *BuildResourceIdFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "build_resource_id"
}

func (b *BuildResourceIdFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "parent_id",
				Description:         "The parent ID of the Azure resource.",
				MarkdownDescription: "The parent ID of the Azure resource.",
			},
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
				Name:                "name",
				Description:         "The name of the Azure resource.",
				MarkdownDescription: "The name of the Azure resource.",
			},
		},
		Return:              function.StringReturn{},
		Summary:             "Builds a generic resource ID.",
		Description:         "This function constructs an Azure resource ID given the parent ID, resource type, and resource name. It is useful for creating resource IDs for top-level and nested resources within a specific scope.",
		MarkdownDescription: "This function constructs an Azure resource ID given the parent ID, resource type, and resource name. It is useful for creating resource IDs for top-level and nested resources within a specific scope.",
	}
}

func (b *BuildResourceIdFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var resourceTypeParam types.String
	var parentId types.String
	var name types.String

	if response.Error = request.Arguments.Get(ctx, &parentId, &resourceTypeParam, &name); response.Error != nil {
		return
	}

	resourceType := utils.TryAppendDefaultApiVersion(resourceTypeParam.ValueString())

	resourceID, err := parse.NewResourceIDSkipScopeValidation(name.ValueString(), parentId.ValueString(), resourceType)
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	response.Error = response.Result.Set(ctx, types.StringValue(resourceID.AzureResourceId))
}

var _ function.Function = &BuildResourceIdFunction{}
