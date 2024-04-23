package myvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type stringIsUUID struct{}

func (v stringIsUUID) Description(ctx context.Context) string {
	return "validate this in UUID format"
}

func (v stringIsUUID) MarkdownDescription(ctx context.Context) string {
	return "validate this in UUID format"
}

func (_ stringIsUUID) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := validation.IsUUID(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid UUID string",
				err.Error())
		}
	}
}

func StringIsUUID() validator.String {
	return stringIsUUID{}
}
