package myvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type stringIsJSON struct{}

func (v stringIsJSON) Description(ctx context.Context) string {
	return "validate this in JSON format"
}

func (v stringIsJSON) MarkdownDescription(ctx context.Context) string {
	return "validate this in JSON format"
}

func (_ stringIsJSON) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := validation.StringIsJSON(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid JSON string",
				err.Error())
		}
	}
}

func StringIsJSON() stringIsJSON {
	return stringIsJSON{}
}
