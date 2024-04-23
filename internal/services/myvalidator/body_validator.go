package myvalidator

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bodyValidator struct{}

func (v bodyValidator) Description(ctx context.Context) string {
	return "validate the identity block"
}

func (v bodyValidator) MarkdownDescription(ctx context.Context) string {
	return "validate the identity block"
}

func (_ bodyValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	raw := req.ConfigValue

	if raw.IsUnknown() || raw.IsNull() || raw.UnderlyingValue() == nil || raw.IsUnderlyingValueNull() || raw.IsUnderlyingValueUnknown() {
		return
	}

	if value, ok := raw.UnderlyingValue().(types.String); ok {
		var out interface{}
		err := json.Unmarshal([]byte(value.ValueString()), &out)
		if err != nil {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON", err.Error())
		}
	}
}

func BodyValidator() validator.Dynamic {
	return bodyValidator{}
}
