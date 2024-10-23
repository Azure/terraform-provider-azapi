package myvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dynamicIsNotStringValidator struct{}

func (v dynamicIsNotStringValidator) Description(ctx context.Context) string {
	return "validate the dynamic value is not a string"
}

func (v dynamicIsNotStringValidator) MarkdownDescription(ctx context.Context) string {
	return "validate the dynamic value is not a string"
}

func (_ dynamicIsNotStringValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	raw := req.ConfigValue

	if raw.IsUnknown() || raw.IsNull() || raw.UnderlyingValue() == nil || raw.IsUnderlyingValueNull() || raw.IsUnderlyingValueUnknown() {
		return
	}

	if _, ok := raw.UnderlyingValue().(types.String); ok {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Type", "The value must not be a string. From 2.0, the value only accepts an HCL object. Please refer to the documentation for more information: https://registry.terraform.io/providers/Azure/azapi/latest/docs/guides/2.0-upgrade-guide#dynamic-properties-support")
	}
}

func DynamicIsNotStringValidator() validator.Dynamic {
	return dynamicIsNotStringValidator{}
}
