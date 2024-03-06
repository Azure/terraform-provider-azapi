package myvalidator

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsUserAssignedIdentityID struct{}

func (v stringIsUserAssignedIdentityID) Description(ctx context.Context) string {
	return "validate this in user assigned identity resource ID format"
}

func (v stringIsUserAssignedIdentityID) MarkdownDescription(ctx context.Context) string {
	return "validate this in user assigned identity resource ID format"
}

func (_ stringIsUserAssignedIdentityID) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := validate.UserAssignedIdentityID(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Resource ID",
				err.Error())
		}
	}
}

func StringIsUserAssignedIdentityID() stringIsUserAssignedIdentityID {
	return stringIsUserAssignedIdentityID{}
}
