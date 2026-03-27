package myvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringIsValidJMESPath_ValidateString(t *testing.T) {
	v := stringIsValidJMESPath{}

	t.Run("valid simple expression", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("properties.status"),
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

	t.Run("valid comparison expression", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("properties.status == 'Running'"),
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

	t.Run("valid complex expression", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("properties.replicaSets[?serviceStatus != 'Running'] | length(@) == `0`"),
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

	t.Run("invalid expression", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("[invalid..expression"),
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

	t.Run("null value", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringNull(),
			Path:        path.Empty(),
		}
		resp := &validator.StringResponse{
			Diagnostics: diag.Diagnostics{},
		}

		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("Expected no errors for null value, but got: %v", resp.Diagnostics)
		}
	})

	t.Run("unknown value", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: basetypes.NewStringUnknown(),
			Path:        path.Empty(),
		}
		resp := &validator.StringResponse{
			Diagnostics: diag.Diagnostics{},
		}

		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("Expected no errors for unknown value, but got: %v", resp.Diagnostics)
		}
	})
}
