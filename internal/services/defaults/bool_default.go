package defaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type boolDefault struct {
	defaultValue bool
}

func (s boolDefault) Description(ctx context.Context) string {
	return "Return the default value"
}

func (s boolDefault) MarkdownDescription(ctx context.Context) string {
	return "Return the default value"
}

func (s boolDefault) DefaultBool(ctx context.Context, request defaults.BoolRequest, response *defaults.BoolResponse) {
	response.PlanValue = basetypes.NewBoolValue(s.defaultValue)
}

func BoolDefault(defaultValue bool) defaults.Bool {
	return boolDefault{defaultValue: defaultValue}
}
