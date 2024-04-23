package myvalidator

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type identityValidator struct{}

func (v identityValidator) Description(ctx context.Context) string {
	return "validate the identity block"
}

func (v identityValidator) MarkdownDescription(ctx context.Context) string {
	return "validate the identity block"
}

func (_ identityValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	value := req.ConfigValue

	if value.IsUnknown() || value.IsNull() {
		return
	}

	var model identity.Model
	if resp.Diagnostics.Append(value.As(ctx, &model, basetypes.ObjectAsOptions{})...); resp.Diagnostics.HasError() {
		return
	}

	if _, err := identity.ExpandIdentity(model); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid `identity` block", err.Error())
	}
}

func IdentityValidator() validator.Object {
	return identityValidator{}
}
