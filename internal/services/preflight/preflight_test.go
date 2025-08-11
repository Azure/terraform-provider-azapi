package preflight

import (
	"regexp"
	"testing"

	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
)

func Test_ParentIdPlaceholder(t *testing.T) {
	testcases := []struct {
		ResourceDef    *aztypes.ResourceType
		SubscriptionId string
		Expected       string
		ExpectedReg    string
		ExpectedErr    bool
	}{
		{
			ResourceDef:    nil,
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "",
			ExpectedErr:    true,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				ScopeTypes: nil,
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "",
			ExpectedErr:    true,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Billing/billingAccounts@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.Tenant},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "/",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Billing/billingAccounts@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.ManagementGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/providers/Microsoft.Management/managementGroups/.+",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Resources/resourceGroups@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.Subscription},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "/subscriptions/00000000-0000-0000-0000-000000000000",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Network/virtualNetworks@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.ResourceGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/.+",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Billing/billingAccounts@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.Tenant, aztypes.ManagementGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "",
			ExpectedErr:    true,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Billing/billingAccounts@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.Extension},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "",
			ExpectedErr:    true,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Network/virtualNetworks/subnets@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.ResourceGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/.+/providers/Microsoft.Network/virtualNetworks/.+",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Network/virtualNetworks/subnets/networkSecurityGroups@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.ResourceGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/.+/providers/Microsoft.Network/virtualNetworks/.+/subnets/.+",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				Name:       "Microsoft.Network/virtualNetworks/subnets/networkSecurityGroups@2000-01-01",
				ScopeTypes: []aztypes.ScopeType{aztypes.Subscription},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/virtualNetworks/.+/subnets/.+",
			ExpectedErr:    false,
		},
	}

	for _, testcase := range testcases {
		actual, err := ParentIdPlaceholder(testcase.ResourceDef, testcase.SubscriptionId)
		if testcase.ExpectedErr {
			if err == nil {
				t.Errorf("Expected error, but got nil")
			}
		} else {
			if testcase.Expected != "" && actual != testcase.Expected {
				t.Errorf("Expected %s, but got %s", testcase.Expected, actual)
			}
			if testcase.ExpectedReg != "" && !regexp.MustCompile(testcase.ExpectedReg).MatchString(actual) {
				t.Errorf("Expected %s, but got %s", testcase.ExpectedReg, actual)
			}
		}
	}
}

func Test_ScopeID(t *testing.T) {
	testcases := []struct {
		ResourceId      string
		ExpectedScopeId string
	}{
		{
			ResourceId:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySubnet",
			ExpectedScopeId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
		},
		{
			ResourceId:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg/providers/Microsoft.Network/virtualNetworks/vnet",
			ExpectedScopeId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
		},
		{
			ResourceId:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
			ExpectedScopeId: "/subscriptions/00000000-0000-0000-0000-000000000000",
		},
		{
			ResourceId:      "/subscriptions/00000000-0000-0000-0000-000000000000",
			ExpectedScopeId: "/",
		},
		{
			ResourceId:      "/providers/Microsoft.Management/managementGroups/azapifakemg",
			ExpectedScopeId: "/",
		},
		{
			ResourceId:      "/providers/Microsoft.Management/managementGroups/azapifakemg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySubnet",
			ExpectedScopeId: "/providers/Microsoft.Management/managementGroups/azapifakemg",
		},
	}

	for _, testcase := range testcases {
		actual, err := ScopeID(testcase.ResourceId)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if actual != testcase.ExpectedScopeId {
			t.Errorf("Expected %s, but got %s", testcase.ExpectedScopeId, actual)
		}
	}
}

func Test_ResourceAndParentName(t *testing.T) {
	testcases := []struct {
		ResourceId                    string
		ScopeId                       string
		ExpectedResourceAndParentName string
	}{
		{
			ResourceId:                    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySubnet",
			ScopeId:                       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
			ExpectedResourceAndParentName: "vnet/mySubnet",
		},
		{
			ResourceId:                    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg/providers/Microsoft.Network/virtualNetworks/vnet",
			ScopeId:                       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
			ExpectedResourceAndParentName: "vnet",
		},
		{
			ResourceId:                    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
			ScopeId:                       "/subscriptions/00000000-0000-0000-0000-000000000000",
			ExpectedResourceAndParentName: "azapifakerg",
		},
		{
			ResourceId:                    "/subscriptions/00000000-0000-0000-0000-000000000000",
			ScopeId:                       "/",
			ExpectedResourceAndParentName: "00000000-0000-0000-0000-000000000000",
		},
		{
			ResourceId:                    "/providers/Microsoft.Management/managementGroups/azapifakemg",
			ScopeId:                       "/",
			ExpectedResourceAndParentName: "azapifakemg",
		},
		{
			ResourceId:                    "/providers/Microsoft.Management/managementGroups/azapifakemg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySubnet",
			ScopeId:                       "/providers/Microsoft.Management/managementGroups/azapifakemg",
			ExpectedResourceAndParentName: "vnet/mySubnet",
		},
	}

	for _, testcase := range testcases {
		actual, err := ResourceAndParentName(testcase.ResourceId, testcase.ScopeId)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if actual != testcase.ExpectedResourceAndParentName {
			t.Errorf("Expected %s, but got %s", testcase.ExpectedResourceAndParentName, actual)
		}
	}
}
