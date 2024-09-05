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
				ScopeTypes: []aztypes.ScopeType{aztypes.Tenant},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "/",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				ScopeTypes: []aztypes.ScopeType{aztypes.ManagementGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/providers/Microsoft.Management/managementGroups/.+",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				ScopeTypes: []aztypes.ScopeType{aztypes.Subscription},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "/subscriptions/00000000-0000-0000-0000-000000000000",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				ScopeTypes: []aztypes.ScopeType{aztypes.ResourceGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			ExpectedReg:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/.+",
			ExpectedErr:    false,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				ScopeTypes: []aztypes.ScopeType{aztypes.Tenant, aztypes.ManagementGroup},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "",
			ExpectedErr:    true,
		},

		{
			ResourceDef: &aztypes.ResourceType{
				ScopeTypes: []aztypes.ScopeType{aztypes.Extension},
			},
			SubscriptionId: "00000000-0000-0000-0000-000000000000",
			Expected:       "",
			ExpectedErr:    true,
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

func Test_IsSupported(t *testing.T) {
	testcases := []struct {
		ParentId     string
		ResourceType string
		Expected     bool
	}{
		{
			// resource type version is not specified
			ParentId:     "",
			ResourceType: "Microsoft.Network/virtualNetworks",
			Expected:     false,
		},

		{
			ParentId:     "",
			ResourceType: "Microsoft.Network/virtualNetworks@2020-06-01",
			Expected:     true,
		},

		{
			// resource is not top level resource
			ParentId:     "",
			ResourceType: "Microsoft.Network/virtualNetworks/subnets@2020-06-01",
			Expected:     false,
		},

		{
			// unknown deploy scope
			ParentId:     "",
			ResourceType: "Microsoft.Resources/deployments@2020-06-01",
			Expected:     false,
		},

		{
			ParentId:     "",
			ResourceType: "Microsoft.Resources/resourceGroups@2020-06-01",
			Expected:     true,
		},

		{
			ParentId:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/azapifakerg",
			ResourceType: "Microsoft.Network/virtualNetworks@2020-06-01",
			Expected:     true,
		},

		{
			ParentId:     "/subscriptions/00000000-0000-0000-0000-000000000000",
			ResourceType: "Microsoft.Network/virtualNetworks@2020-06-01",
			Expected:     true,
		},

		{
			ParentId:     "/",
			ResourceType: "Microsoft.Network/virtualNetworks@2020-06-01",
			Expected:     true,
		},

		{
			ParentId:     "/providers/Microsoft.Management/managementGroups/azapifakemg",
			ResourceType: "Microsoft.Network/virtualNetworks@2020-06-01",
			Expected:     true,
		},
	}

	for _, testcase := range testcases {
		actual := IsSupported(testcase.ResourceType, testcase.ParentId)
		if actual != testcase.Expected {
			t.Errorf("Expected %v, but got %v", testcase.Expected, actual)
		}
	}
}
