package defaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type dynamicDefault struct {
	defaultValue attr.Value
}

func (s dynamicDefault) Description(ctx context.Context) string {
	return "Return the default value"
}

func (s dynamicDefault) MarkdownDescription(ctx context.Context) string {
	return "Return the default value"
}

func (s dynamicDefault) DefaultDynamic(ctx context.Context, request defaults.DynamicRequest, response *defaults.DynamicResponse) {
	response.PlanValue = basetypes.NewDynamicValue(s.defaultValue)
}

func DynamicDefault(defaultValue attr.Value) defaults.Dynamic {
	return dynamicDefault{defaultValue: defaultValue}
}
