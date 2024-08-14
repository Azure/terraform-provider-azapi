package functions_test

import (
	"context"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/functions"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_BuildResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"build-subscription-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000"),
					types.StringValue("Microsoft.Resources/resourceGroups"),
					types.StringValue("myResourceGroup"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup")),
			},
		},
		"build-resource-group-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"),
					types.StringValue("Microsoft.Automation/automationAccounts"),
					types.StringValue("myAutomationAccount"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Automation/automationAccounts/myAutomationAccount")),
			},
		},
		"build-management-group-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/providers/Microsoft.Management/managementGroups/mg1"),
					types.StringValue("Microsoft.Authorization/privateLinkAssociations"),
					types.StringValue("MyPrivateLinkAssociation"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Authorization/privateLinkAssociations/MyPrivateLinkAssociation")),
			},
		},
		"build-tenant-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/"),
					types.StringValue("Microsoft.Billing/billingAccounts"),
					types.StringValue("MyInvoice"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Billing/billingAccounts/MyInvoice")),
			},
		},
		"build-extension-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM"),
					types.StringValue("Microsoft.Chaos/targets"),
					types.StringValue("MyTarget"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM/providers/Microsoft.Chaos/targets/MyTarget")),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := function.RunResponse{
				Result: function.NewResultData(types.StringUnknown()),
			}

			buildResourceIdFunction := functions.BuildResourceIdFunction{}
			buildResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
