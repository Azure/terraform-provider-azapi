package services

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

const (
	casingTestName     = "plan1"
	casingTestParentID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
	casingTestType     = "Microsoft.Web/serverfarms@2025-03-01"

	// State migrated from azurerm carries the casing Azure returns ("serverFarms"),
	// while the provider rebuilds the id from the configured type ("serverfarms").
	casingTestMigratedID   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Web/serverFarms/plan1"
	casingTestRecomputedID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Web/serverfarms/plan1"
)

// TestPreserveResourceIDCasingAzapiResource covers the azapi_resource create/update
// path, where the id is rebuilt with parse.NewResourceID and assigned to plan.ID.
func TestPreserveResourceIDCasingAzapiResource(t *testing.T) {
	id, err := parse.NewResourceID(casingTestName, casingTestParentID, casingTestType)
	if err != nil {
		t.Fatalf("building resource id: %+v", err)
	}

	// Without this the assertions below would pass vacuously.
	if id.ID() != casingTestRecomputedID {
		t.Fatalf("parse.NewResourceID() = %q, expected %q", id.ID(), casingTestRecomputedID)
	}

	if got := preserveCasing(casingTestMigratedID, id.ID(), false); got != casingTestRecomputedID {
		t.Errorf("feature disabled: got %q, expected %q", got, casingTestRecomputedID)
	}

	if got := preserveCasing(casingTestMigratedID, id.ID(), true); got != casingTestMigratedID {
		t.Errorf("feature enabled: got %q, expected %q", got, casingTestMigratedID)
	}
}

// TestPreserveResourceIDCasingAzapiUpdateResourceNameParentID covers the
// azapi_update_resource create/update path when name and parent_id are configured.
// The last type segment comes from type, so both id and resource_id are rebuilt.
func TestPreserveResourceIDCasingAzapiUpdateResourceNameParentID(t *testing.T) {
	id, err := parse.NewResourceID(casingTestName, casingTestParentID, casingTestType)
	if err != nil {
		t.Fatalf("building resource id: %+v", err)
	}

	if id.ID() != casingTestRecomputedID || id.AzureResourceId != casingTestRecomputedID {
		t.Fatalf("parse.NewResourceID() = (%q, %q), expected both to be %q", id.ID(), id.AzureResourceId, casingTestRecomputedID)
	}

	testCases := []struct {
		name      string
		candidate string
		enabled   bool
		expected  string
	}{
		{"id disabled", id.ID(), false, casingTestRecomputedID},
		{"id enabled", id.ID(), true, casingTestMigratedID},
		{"resource_id disabled", id.AzureResourceId, false, casingTestRecomputedID},
		{"resource_id enabled", id.AzureResourceId, true, casingTestMigratedID},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := preserveCasing(casingTestMigratedID, tc.candidate, tc.enabled); got != tc.expected {
				t.Fatalf("got %q, expected %q", got, tc.expected)
			}
		})
	}
}

// TestPreserveResourceIDCasingAzapiUpdateResourceResourceID covers the
// azapi_update_resource path when resource_id is configured. parse.ResourceIDWithResourceType
// currently echoes the supplied casing, so this asserts the guarantee the feature
// makes rather than that behaviour, and fails if a future parser change reintroduces #1120.
func TestPreserveResourceIDCasingAzapiUpdateResourceResourceID(t *testing.T) {
	id, err := parse.ResourceIDWithResourceType(casingTestMigratedID, casingTestType)
	if err != nil {
		t.Fatalf("parsing resource id: %+v", err)
	}

	if got := preserveCasing(casingTestMigratedID, id.ID(), true); got != casingTestMigratedID {
		t.Errorf("feature enabled: got %q, expected %q", got, casingTestMigratedID)
	}

	if got := preserveCasing(casingTestMigratedID, id.AzureResourceId, true); got != casingTestMigratedID {
		t.Errorf("feature enabled: got %q, expected %q", got, casingTestMigratedID)
	}
}
