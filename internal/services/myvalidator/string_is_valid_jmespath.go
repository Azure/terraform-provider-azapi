package myvalidator

import (
	"context"

	jmes "github.com/jmespath/go-jmespath"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsValidJMESPath struct{}

func (v stringIsValidJMESPath) Description(ctx context.Context) string {
	return "validates that the string compiles as a valid JMESPath expression"
}

func (v stringIsValidJMESPath) MarkdownDescription(ctx context.Context) string {
	return "validates that the string compiles as a valid JMESPath expression"
}

func (stringIsValidJMESPath) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	strval := str.ValueString()
	if _, err := jmes.Compile(strval); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JMESPath expression",
			err.Error(),
		)
	}
}

func StringIsValidJMESPath() validator.String {
	return stringIsValidJMESPath{}
}
