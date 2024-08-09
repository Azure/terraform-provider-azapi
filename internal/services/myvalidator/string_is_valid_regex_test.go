package myvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringIsValidRegex_ValidateString(t *testing.T) {
	v := stringIsValidRegex{}

	t.Run("valid regex", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("^[a-zA-Z]+$"),
			Path:        path.Empty(),
		}
		resp := &validator.StringResponse{
			Diagnostics: diag.Diagnostics{},
		}

		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("Expected no errors, but got: %v", resp.Diagnostics)
		}
	})

	t.Run("invalid regex", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("[a-zA-Z+$"),
			Path:        path.Empty(),
		}
		resp := &validator.StringResponse{
			Diagnostics: diag.Diagnostics{},
		}

		v.ValidateString(context.Background(), req, resp)

		if !resp.Diagnostics.HasError() {
			t.Errorf("Expected errors, but got none")
		}
	})
}
