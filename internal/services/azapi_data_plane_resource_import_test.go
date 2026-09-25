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

	t.Run("foundry dataset version", func(t *testing.T) {
		resourceID, resourceType, err := parseDataPlaneImportID(
			"contoso.services.ai.azure.com/api/projects/example/datasets/training/versions/1|Microsoft.Foundry/datasets/versions@2025-05-01",
		)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		id, err := parse.DataPlaneResourceIDWithResourceType(resourceID, resourceType)
		if err != nil {
			t.Fatalf("expected dataset ID to parse, got: %v", err)
		}
		if id.Name != "1" {
			t.Fatalf("unexpected dataset version: %q", id.Name)
		}
		if id.ParentId != "contoso.services.ai.azure.com/api/projects/example/datasets/training" {
			t.Fatalf("unexpected dataset parent: %q", id.ParentId)
		}
	})
}
