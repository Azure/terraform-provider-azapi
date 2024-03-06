package tags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ExpandTags(input types.Map) map[string]string {
	output := make(map[string]string)
	if diags := input.ElementsAs(context.Background(), &output, false); diags.HasError() {
		return nil
	}
	return output
}

func FlattenTags(input interface{}) types.Map {
	if input == nil {
		return basetypes.NewMapNull(types.StringType)
	}
	tagsMap := make(map[string]attr.Value)
	switch inputMap := input.(type) {
	case map[string]interface{}:
		for k, v := range inputMap {
			tagsMap[k] = basetypes.NewStringValue(v.(string))
		}
	case map[string]string:
		for k, v := range inputMap {
			tagsMap[k] = basetypes.NewStringValue(v)
		}
	}
	return basetypes.NewMapValueMust(types.StringType, tagsMap)
}

func Validator() validator.Map {
	return tagsValidator{}
}

type tagsValidator struct{}

func (v tagsValidator) ValidateMap(ctx context.Context, request validator.MapRequest, response *validator.MapResponse) {
	input := request.ConfigValue

	if input.IsUnknown() || input.IsNull() {
		return
	}

	tagsMap := make(map[string]string)
	if response.Diagnostics.Append(input.ElementsAs(context.Background(), &tagsMap, false)...); response.Diagnostics.HasError() {
		return
	}

	if len(tagsMap) > 50 {
		response.Diagnostics.AddError("Invalid tags", fmt.Errorf("a maximum of 50 tags can be applied to each ARM resource").Error())
	}

	for k, v := range tagsMap {
		if len(k) > 512 {
			response.Diagnostics.AddError("Invalid tags", fmt.Errorf("the maximum length for a tag key is 512 characters: %q is %d characters", k, len(k)).Error())
		}

		if len(v) > 256 {
			response.Diagnostics.AddError("Invalid tags", fmt.Errorf("the maximum length for a tag value is 256 characters: the value for %q is %d characters", k, len(v)).Error())
		}
	}

}

func (v tagsValidator) Description(ctx context.Context) string {
	return "validate the tags"
}

func (v tagsValidator) MarkdownDescription(ctx context.Context) string {
	return "validate the tags"
}
