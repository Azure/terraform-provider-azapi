package myvalidator

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsResourceID struct{}

func (v stringIsResourceID) Description(ctx context.Context) string {
	return "validate this in resource ID format"
}

func (v stringIsResourceID) MarkdownDescription(ctx context.Context) string {
	return "validate this in resource ID format"
}

func (_ stringIsResourceID) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := validate.ResourceID(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Resource ID",
				err.Error())
		}
	}
}

func StringIsResourceID() stringIsResourceID {
	return stringIsResourceID{}
}
