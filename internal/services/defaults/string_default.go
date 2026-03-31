package defaults

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type stringDefault struct {
	defaultValue string
}

func (s stringDefault) Description(ctx context.Context) string {
	return fmt.Sprintf("defaults to %s", s.defaultValue)
}

func (s stringDefault) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("defaults to `%s`", s.defaultValue)
}

func (s stringDefault) DefaultString(ctx context.Context, request defaults.StringRequest, response *defaults.StringResponse) {
	response.PlanValue = basetypes.NewStringValue(s.defaultValue)
}

func StringDefault(defaultValue string) defaults.String {
	return stringDefault{defaultValue: defaultValue}
}
