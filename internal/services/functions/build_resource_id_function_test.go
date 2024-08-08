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
					types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
					types.StringValue("myResourceGroup"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.BuildResourceIdResultAttrTypes, map[string]attr.Value{
					"resource_id": types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"),
				})),
			},
		},
		"build-subscription-scope-resource-id-invalid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Invalid"),
					types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
					types.StringValue("myResourceGroup"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`parent_id is invalid`: expect ID of resource whose scope is [Subscription], but got scope Unknown"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"build-resource-group-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"),
					types.StringValue("Microsoft.Automation/automationAccounts@2021-06-22"),
					types.StringValue("myAutomationAccount"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.BuildResourceIdResultAttrTypes, map[string]attr.Value{
					"resource_id": types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Automation/automationAccounts/myAutomationAccount"),
				})),
			},
		},
		"build-resource-group-scope-resource-id-invalid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Invalid"),
					types.StringValue("Microsoft.Automation/automationAccounts@2021-06-22"),
					types.StringValue("myAutomationAccount"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`parent_id is invalid`: expect ID of resource whose scope is [ResourceGroup], but got scope Unknown"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"build-management-group-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/providers/Microsoft.Management/managementGroups/mg1"),
					types.StringValue("Microsoft.Authorization/privateLinkAssociations@2020-05-01"),
					types.StringValue("MyPrivateLinkAssociation"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.BuildResourceIdResultAttrTypes, map[string]attr.Value{
					"resource_id": types.StringValue("/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Authorization/privateLinkAssociations/MyPrivateLinkAssociation"),
				})),
			},
		},
		"build-management-group-scope-resource-id-invalid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Invalid"),
					types.StringValue("Microsoft.Authorization/privateLinkAssociations@2020-05-01"),
					types.StringValue("MyPrivateLinkAssociation"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`parent_id is invalid`: expect ID of resource whose scope is [ManagementGroup], but got scope Unknown"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"build-tenant-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/"),
					types.StringValue("Microsoft.Billing/billingAccounts@2018-11-01-preview"),
					types.StringValue("MyInvoice"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.BuildResourceIdResultAttrTypes, map[string]attr.Value{
					"resource_id": types.StringValue("/providers/Microsoft.Billing/billingAccounts/MyInvoice"),
				})),
			},
		},
		"build-tenant-scope-resource-id-invalid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Invalid"),
					types.StringValue("Microsoft.Billing/billingAccounts@2018-11-01-preview"),
					types.StringValue("MyInvoice"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`parent_id is invalid`: expect ID of resource whose scope is [Tenant], but got scope Unknown"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"build-extension-scope-resource-id-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM"),
					types.StringValue("Microsoft.Chaos/targets@2024-01-01"),
					types.StringValue("MyTarget"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.BuildResourceIdResultAttrTypes, map[string]attr.Value{
					"resource_id": types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM/providers/Microsoft.Chaos/targets/MyTarget"),
				})),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := function.RunResponse{
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			}

			buildResourceIdFunction := functions.BuildResourceIdFunction{}
			buildResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
