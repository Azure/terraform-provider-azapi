package services

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateDataPlaneResourceName(t *testing.T) {
	t.Run("agent requires name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/agents@v1"),
			Name: types.StringNull(),
		}

		err := validateDataPlaneResourceName(config)
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "must be set") {
			t.Fatalf("expected must-be-set error, got: %v", err)
		}
	})

	t.Run("agent accepts explicit name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/agents@v1"),
			Name: types.StringValue("terraform-agent"),
		}

		if err := validateDataPlaneResourceName(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("non-assistant requires name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.KeyVault/vaults/secrets@7.4"),
			Name: types.StringNull(),
		}

		err := validateDataPlaneResourceName(config)
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "must be set") {
			t.Fatalf("expected must-be-set error, got: %v", err)
		}
	})

	t.Run("non-assistant accepts name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.KeyVault/vaults/secrets@7.4"),
			Name: types.StringValue("secret-name"),
		}

		if err := validateDataPlaneResourceName(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})
}
