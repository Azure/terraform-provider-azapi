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

	t.Run("table entity does not take a name (composite key lives in parent_id)", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"),
			Name: types.StringNull(),
		}

		if err := validateDataPlaneResourceName(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("table entity rejects name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"),
			Name: types.StringValue("unexpected"),
		}

		err := validateDataPlaneResourceName(config)
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "does not have a name") {
			t.Fatalf("expected no-name error, got: %v", err)
		}
	})
}

func TestValidateDataPlaneResourceWritable(t *testing.T) {
	t.Run("read-only entities collection is rejected", func(t *testing.T) {
		err := validateDataPlaneResourceWritable("Microsoft.Storage/storageAccounts/tableServices/tables/entitiesCollection@2026-04-06")
		if err == nil {
			t.Fatalf("expected validation error")
		}
		if !strings.Contains(err.Error(), "does not support create/update/delete") {
			t.Fatalf("expected read-only error, got: %v", err)
		}
	})

	t.Run("writable table type is accepted", func(t *testing.T) {
		if err := validateDataPlaneResourceWritable("Microsoft.Storage/storageAccounts/tableServices/tables@2026-04-06"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("type without customization is accepted", func(t *testing.T) {
		if err := validateDataPlaneResourceWritable("Microsoft.Foundry/agents@v1"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})
}
