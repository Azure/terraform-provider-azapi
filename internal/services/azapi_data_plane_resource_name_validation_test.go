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

	t.Run("foundry evaluation generates name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/evaluation/versions@2025-05-01"),
			Name: types.StringNull(),
		}

		if err := validateDataPlaneResourceName(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("foundry evaluation rejects configured name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/evaluation/versions@2025-05-01"),
			Name: types.StringValue("evaluation-name"),
		}

		err := validateDataPlaneResourceName(config)
		if err == nil {
			t.Fatal("expected validation error")
		}
		if !strings.Contains(err.Error(), "should not be set") {
			t.Fatalf("expected generated-name error, got: %v", err)
		}
	})

	t.Run("foundry evaluation run generates name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/evaluation/runs@2025-05-01"),
			Name: types.StringNull(),
		}

		if err := validateDataPlaneResourceName(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("foundry evaluation run rejects configured name", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type: types.StringValue("Microsoft.Foundry/evaluation/runs@2025-05-01"),
			Name: types.StringValue("run-name"),
		}

		err := validateDataPlaneResourceName(config)
		if err == nil {
			t.Fatal("expected validation error")
		}
		if !strings.Contains(err.Error(), "should not be set") {
			t.Fatalf("expected generated-name error, got: %v", err)
		}
	})
}

func TestValidateDataPlaneResourceLookupName(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		resourceName types.String
		wantError    bool
	}{
		{
			name:         "evaluation ID is accepted for lookup",
			resourceType: "Microsoft.Foundry/evaluation/versions@2025-05-01",
			resourceName: types.StringValue("eval_123"),
		},
		{
			name:         "evaluation ID is required for lookup",
			resourceType: "Microsoft.Foundry/evaluation/versions@2025-05-01",
			resourceName: types.StringNull(),
			wantError:    true,
		},
		{
			name:         "run ID is accepted for lookup",
			resourceType: "Microsoft.Foundry/evaluation/runs@2025-05-01",
			resourceName: types.StringValue("run_456"),
		},
		{
			name:         "run ID is required for lookup",
			resourceType: "Microsoft.Foundry/evaluation/runs@2025-05-01",
			resourceName: types.StringNull(),
			wantError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateDataPlaneResourceLookupName(&DataPlaneResourceModel{
				Type: types.StringValue(test.resourceType),
				Name: test.resourceName,
			})
			if test.wantError {
				if err == nil {
					t.Fatal("expected validation error")
				}
				if !strings.Contains(err.Error(), `argument "name"`) {
					t.Fatalf("expected name validation error, got: %v", err)
				}
			} else if err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
		})
	}
}

func TestValidateDataPlaneResourceEvaluationID(t *testing.T) {
	t.Run("evaluation run requires evaluation ID", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type:         types.StringValue("Microsoft.Foundry/evaluation/runs@2025-05-01"),
			EvaluationID: types.StringNull(),
		}

		err := validateDataPlaneResourceEvaluationID(config)
		if err == nil {
			t.Fatal("expected validation error")
		}
		if !strings.Contains(err.Error(), "evaluation_id") {
			t.Fatalf("expected evaluation_id error, got: %v", err)
		}
	})

	t.Run("evaluation run accepts evaluation ID", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type:         types.StringValue("Microsoft.Foundry/evaluation/runs@2025-05-01"),
			EvaluationID: types.StringValue("eval_123"),
		}

		if err := validateDataPlaneResourceEvaluationID(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("evaluation version ignores evaluation ID", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type:         types.StringValue("Microsoft.Foundry/evaluation/versions@2025-05-01"),
			EvaluationID: types.StringNull(),
		}

		if err := validateDataPlaneResourceEvaluationID(config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("evaluation version rejects evaluation ID", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type:         types.StringValue("Microsoft.Foundry/evaluation/versions@2025-05-01"),
			EvaluationID: types.StringValue("eval_123"),
		}

		err := validateDataPlaneResourceEvaluationID(config)
		if err == nil {
			t.Fatal("expected validation error")
		}
		if !strings.Contains(err.Error(), "should not be set") {
			t.Fatalf("expected unsupported evaluation ID error, got: %v", err)
		}
	})

	t.Run("non-Foundry resource rejects evaluation ID", func(t *testing.T) {
		config := &DataPlaneResourceModel{
			Type:         types.StringValue("Microsoft.KeyVault/vaults/secrets@7.4"),
			EvaluationID: types.StringValue("eval_123"),
		}

		err := validateDataPlaneResourceEvaluationID(config)
		if err == nil {
			t.Fatal("expected validation error")
		}
		if !strings.Contains(err.Error(), "should not be set") {
			t.Fatalf("expected unsupported evaluation ID error, got: %v", err)
		}
	})
}
