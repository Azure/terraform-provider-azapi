package defaults

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type dynamicDefault struct {
	defaultValue attr.Value
}

func (s dynamicDefault) Description(ctx context.Context) string {
	return fmt.Sprintf("defaults to %s", s.defaultValue.String())
}

func (s dynamicDefault) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("defaults to `%s`", s.defaultValue.String())
}

func (s dynamicDefault) DefaultDynamic(ctx context.Context, request defaults.DynamicRequest, response *defaults.DynamicResponse) {
	response.PlanValue = basetypes.NewDynamicValue(s.defaultValue)
}

func DynamicDefault(defaultValue attr.Value) defaults.Dynamic {
	return dynamicDefault{defaultValue: defaultValue}
}
