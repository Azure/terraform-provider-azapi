package utils_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
)

func Test_GetParentId(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "",
			Output: "",
		},
		{
			Input:  "/",
			Output: "",
		},
		{
			Input:  "/providers/Microsoft.Billing/billingAccounts/myAccount",
			Output: "",
		},
		{
			Input:  "/providers/Microsoft.Billing/billingAccounts/myAccount/billingProfiles/myProfile",
			Output: "/providers/Microsoft.Billing/billingAccounts/myAccount",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012",
			Output: "",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/policydefinitions/myDef",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/featureProviders/myProvider/subscriptionFeatureRegistrations/myFeature",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/featureProviders/myProvider",
		},
		{
			Input:  "/providers/Microsoft.Management/managementGroups/myMgmtGroup",
			Output: "",
		},
		{
			Input:  "/providers/Microsoft.Management/managementGroups/myMgmtGroup/providers/Microsoft.CostManagement/externalSubscriptions/mySub",
			Output: "/providers/Microsoft.Management/managementGroups/myMgmtGroup",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/sg1",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/sg1/providers/Microsoft.Insights/metrics/m1",
			Output: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/sg1",
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
			Input:  "",
			Output: "Tenant",
		},
		{
			Input:  "/",
			Output: "Tenant",
		},
		{
			Input:  "/providers/Microsoft.Billing/billingAccounts/myAccount",
			Output: "Microsoft.Billing/billingAccounts",
		},
		{
			Input:  "/providers/Microsoft.Billing/billingAccounts/myAccount/billingProfiles/myProfile",
			Output: "Microsoft.Billing/billingAccounts/billingProfiles",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012",
			Output: "Subscription",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/policydefinitions/myDef",
			Output: "Microsoft.Authorization/policydefinitions",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/featureProviders/myProvider/subscriptionFeatureRegistrations/myFeature",
			Output: "Microsoft.Features/featureProviders/subscriptionFeatureRegistrations",
		},
		{
			Input:  "/providers/Microsoft.Management/managementGroups/myMgmtGroup",
			Output: "Microsoft.Management/managementGroups",
		},
		{
			Input:  "/providers/Microsoft.Management/managementGroups/myMgmtGroup/providers/Microsoft.CostManagement/externalSubscriptions/mySub",
			Output: "Microsoft.CostManagement/externalSubscriptions",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1",
			Output: "Microsoft.Resources/resourceGroups",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/sg1",
			Output: "Microsoft.Network/networkSecurityGroups",
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/sg1/providers/Microsoft.Insights/metrics/m1",
			Output: "Microsoft.Insights/metrics",
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
			Output: "",
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

func Test_GetScopeType(t *testing.T) {
	cases := []struct {
		Input  string
		Output types.ScopeType
	}{
		{
			Input:  "",
			Output: types.Tenant,
		},
		{
			Input:  "/",
			Output: types.Tenant,
		},
		{
			Input:  "/providers/Microsoft.Billing/billingAccounts/myAccount",
			Output: types.Tenant,
		},
		{
			Input:  "/providers/Microsoft.Billing/billingAccounts/myAccount/billingProfiles/myProfile",
			Output: types.Tenant,
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012",
			Output: types.Subscription,
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/policydefinitions/myDef",
			Output: types.Subscription,
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Features/featureProviders/myProvider/subscriptionFeatureRegistrations/myFeature",
			Output: types.Subscription,
		},
		{
			Input:  "/providers/Microsoft.Management/managementGroups/myMgmtGroup",
			Output: types.ManagementGroup,
		},
		{
			Input:  "/providers/Microsoft.Management/managementGroups/myMgmtGroup/providers/Microsoft.CostManagement/externalSubscriptions/mySub",
			Output: types.ManagementGroup,
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1",
			Output: types.ResourceGroup,
		},
		{
			Input:  "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/sg1",
			Output: types.ResourceGroup,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		output := utils.GetScopeType(tc.Input)

		if tc.Output != output {
			t.Fatalf("Expected %v but got %v", tc.Output, output)
		}
	}
}
