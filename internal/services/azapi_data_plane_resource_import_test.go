package services

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func TestParseDataPlaneImportID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		resourceID, resourceType, err := parseDataPlaneImportID("host/api/projects/project/agents/agent|Microsoft.Foundry/agents@v1")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
		if resourceID != "host/api/projects/project/agents/agent" {
			t.Fatalf("unexpected resource ID: %q", resourceID)
		}
		if resourceType != "Microsoft.Foundry/agents@v1" {
			t.Fatalf("unexpected resource type: %q", resourceType)
		}
	})

	t.Run("missingType", func(t *testing.T) {
		_, _, err := parseDataPlaneImportID("host/api/projects/project/agents/agent")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("invalidType", func(t *testing.T) {
		_, _, err := parseDataPlaneImportID("host/api/projects/project/agents/agent|Microsoft.Foundry/agents")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("foundry evaluation run", func(t *testing.T) {
		resourceID, resourceType, err := parseDataPlaneImportID(
			"host/api/projects/project/openai/v1/evals/eval_123/runs/run_456|Microsoft.Foundry/evaluation/runs@2025-05-01",
		)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		id, err := parse.DataPlaneResourceIDWithResourceType(resourceID, resourceType)
		if err != nil {
			t.Fatalf("expected evaluation run ID to parse, got: %v", err)
		}
		if id.Name != "run_456" {
			t.Fatalf("unexpected run ID: %q", id.Name)
		}
		if id.EvaluationId != "eval_123" {
			t.Fatalf("unexpected evaluation ID: %q", id.EvaluationId)
		}
		if id.ParentId != "host/api/projects/project" {
			t.Fatalf("unexpected parent: %q", id.ParentId)
		}
	})
}
