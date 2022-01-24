package utils_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
)

func Test_GetParentId(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1/pools/pool1",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1",
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		output := utils.GetParentId(tc.Input)

		if tc.Output != output {
			t.Fatalf("Expected %s but got %s", tc.Output, output)
		}
	}
}

func Test_GetResourceType(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1",
			Output: "Microsoft.Resources/resourceGroups",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1",
			Output: "Microsoft.EventHub/clusters",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1/pools/pool1",
			Output: "Microsoft.EventHub/clusters/pools",
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		output := utils.GetResourceType(tc.Input)

		if tc.Output != output {
			t.Fatalf("Expected %s but got %s", tc.Output, output)
		}
	}
}

func Test_GetParentType(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "Microsoft.EventHub/clusters",
			Output: "Microsoft.Resources/resourceGroups",
		},
		{
			Input:  "Microsoft.EventHub/clusters/pools",
			Output: "Microsoft.EventHub/clusters",
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		output := utils.GetParentType(tc.Input)

		if tc.Output != output {
			t.Fatalf("Expected %s but got %s", tc.Output, output)
		}
	}
}
