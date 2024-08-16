package myvalidator

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsValidRegex struct{}

func (v stringIsValidRegex) Description(ctx context.Context) string {
	return "validates that the string compiles as a valid Go regular expression"
}

func (v stringIsValidRegex) MarkdownDescription(ctx context.Context) string {
	return "validates that the string compiles as a valid Go regular expression"
}

func (stringIsValidRegex) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	strval := str.ValueString()
	if _, err := regexp.Compile(strval); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid regular expression",
			err.Error(),
		)
	}
}

func StringIsValidRegex() validator.String {
	return stringIsValidRegex{}
}
