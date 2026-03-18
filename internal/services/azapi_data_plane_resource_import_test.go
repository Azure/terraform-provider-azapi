package services

import (
	"testing"
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
}
