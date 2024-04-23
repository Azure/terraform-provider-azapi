package myvalidator

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsResourceType struct{}

func (v stringIsResourceType) Description(ctx context.Context) string {
	return "validate this in resource type format"
}

func (v stringIsResourceType) MarkdownDescription(ctx context.Context) string {
	return "validate this in resource type format"
}

func (_ stringIsResourceType) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := validate.ResourceType(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Resource Type",
				err.Error())
		}
	}
}

func StringIsResourceType() validator.String {
	return stringIsResourceType{}
}
