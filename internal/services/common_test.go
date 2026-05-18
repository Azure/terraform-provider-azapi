package services

import "testing"

func TestPreserveCasing(t *testing.T) {
	cases := []struct {
		name      string
		existing  string
		candidate string
		enabled   bool
		want      string
	}{
		{
			name:      "disabled returns candidate even when only casing differs",
			existing:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
			candidate: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverfarms/plan1",
			enabled:   false,
			want:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverfarms/plan1",
		},
		{
			name:      "enabled with empty existing returns candidate",
			existing:  "",
			candidate: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverfarms/plan1",
			enabled:   true,
			want:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverfarms/plan1",
		},
		{
			name:      "enabled and case-insensitively equal returns existing casing",
			existing:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
			candidate: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverfarms/plan1",
			enabled:   true,
			want:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
		},
		{
			name:      "enabled and identical returns existing",
			existing:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
			candidate: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
			enabled:   true,
			want:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
		},
		{
			name:      "enabled but values genuinely differ returns candidate",
			existing:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan1",
			candidate: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan2",
			enabled:   true,
			want:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/serverFarms/plan2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := preserveCasing(tc.existing, tc.candidate, tc.enabled)
			if got != tc.want {
				t.Errorf("preserveCasing(%q, %q, %t) = %q, want %q", tc.existing, tc.candidate, tc.enabled, got, tc.want)
			}
		})
	}
}
