package validate_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
)

func TestResourceID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// a valid tenant Id
			Input: "/",
			Valid: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// a valid subscription Id
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012",
			Valid: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// a valid ResourceGroup Id
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1",
			Valid: true,
		},

		{
			// missing resource type
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/",
			Valid: false,
		},

		{
			// missing value for resource type
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1",
			Valid: true,
		},

		{
			// should not contain api-version
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1?api-version=2020-12-01",
			Valid: false,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.EVENTHUB/CLUSTERS/CLUSTER1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := validate.ResourceID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
