package services

import "testing"

func TestPreserveCasing(t *testing.T) {
	const (
		serverFarmsID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Web/serverFarms/plan1"
		serverfarmsID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Web/serverfarms/plan1"
	)

	testCases := []struct {
		name      string
		existing  string
		candidate string
		enabled   bool
		expected  string
	}{
		{
			name:      "disabled returns candidate",
			existing:  "/foo/Bar",
			candidate: "/foo/bar",
			enabled:   false,
			expected:  "/foo/bar",
		},
		{
			name:      "enabled with empty existing returns candidate",
			existing:  "",
			candidate: "/foo/bar",
			enabled:   true,
			expected:  "/foo/bar",
		},
		{
			name:      "enabled with case-insensitive match keeps existing",
			existing:  "/foo/Bar",
			candidate: "/foo/bar",
			enabled:   true,
			expected:  "/foo/Bar",
		},
		{
			name:      "enabled with identical values keeps existing",
			existing:  "/foo/bar",
			candidate: "/foo/bar",
			enabled:   true,
			expected:  "/foo/bar",
		},
		{
			name:      "enabled with genuinely different values returns candidate",
			existing:  "/foo/bar",
			candidate: "/foo/baz",
			enabled:   true,
			expected:  "/foo/baz",
		},
		{
			// Regression test for #1120: after a state migration the resource id
			// may have been normalised to lower case (serverfarms) while the value
			// the provider recomputes preserves the original casing (serverFarms).
			// With the feature enabled the state casing must be preserved so the
			// azapi_resource identity does not report a spurious change.
			name:      "enabled preserves migrated serverFarms casing (#1120)",
			existing:  serverFarmsID,
			candidate: serverfarmsID,
			enabled:   true,
			expected:  serverFarmsID,
		},
		{
			// Same scenario but with the feature disabled (the default): the
			// recomputed candidate wins so behaviour is unchanged for existing users.
			name:      "disabled does not preserve migrated serverFarms casing (default)",
			existing:  serverFarmsID,
			candidate: serverfarmsID,
			enabled:   false,
			expected:  serverfarmsID,
		},
		{
			// azapi_update_resource.resource_id path: the bare Azure resource id
			// (without the api-version suffix) must also keep its state casing.
			name:      "enabled preserves resource_id path casing",
			existing:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG1/providers/Microsoft.Web/serverFarms/plan1",
			candidate: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Web/serverfarms/plan1",
			enabled:   true,
			expected:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG1/providers/Microsoft.Web/serverFarms/plan1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := preserveCasing(tc.existing, tc.candidate, tc.enabled)
			if actual != tc.expected {
				t.Fatalf("preserveCasing(%q, %q, %t) = %q, expected %q", tc.existing, tc.candidate, tc.enabled, actual, tc.expected)
			}
		})
	}
}
