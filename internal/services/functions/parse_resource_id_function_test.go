package functions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/functions"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"valid-resource-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Network/virtualNetworks"),
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.ParseResourceIdResultAttrTypes, map[string]attr.Value{
					"id":                  types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"),
					"type":                types.StringValue("Microsoft.Network/virtualNetworks"),
					"name":                types.StringValue("vnet1"),
					"parent_id":           types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"),
					"resource_group_name": types.StringValue("rg1"),
					"subscription_id":     types.StringValue("00000000-0000-0000-0000-000000000000"),
					"provider_namespace":  types.StringValue("Microsoft.Network"),
					"parts": types.MapValueMust(types.StringType, map[string]attr.Value{
						"subscriptions":   types.StringValue("00000000-0000-0000-0000-000000000000"),
						"resourceGroups":  types.StringValue("rg1"),
						"providers":       types.StringValue("Microsoft.Network"),
						"virtualNetworks": types.StringValue("vnet1"),
					}),
				})),
			},
		},
		"valid-tenant-resource-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Billing/billingAccounts"),
					types.StringValue("/providers/Microsoft.Billing/billingAccounts/ba1"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.ParseResourceIdResultAttrTypes, map[string]attr.Value{
					"id":                  types.StringValue("/providers/Microsoft.Billing/billingAccounts/ba1"),
					"type":                types.StringValue("Microsoft.Billing/billingAccounts"),
					"name":                types.StringValue("ba1"),
					"parent_id":           types.StringValue("/"),
					"resource_group_name": types.StringValue(""),
					"subscription_id":     types.StringValue(""),
					"provider_namespace":  types.StringValue("Microsoft.Billing"),
					"parts": types.MapValueMust(types.StringType, map[string]attr.Value{
						"billingAccounts": types.StringValue("ba1"),
						"providers":       types.StringValue("Microsoft.Billing"),
					}),
				})),
			},
		},
		"valid-subscription-resource-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Resources/subscriptions"),
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.ParseResourceIdResultAttrTypes, map[string]attr.Value{
					"id":                  types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000"),
					"type":                types.StringValue("Microsoft.Resources/subscriptions"),
					"name":                types.StringValue("00000000-0000-0000-0000-000000000000"),
					"parent_id":           types.StringValue("/"),
					"resource_group_name": types.StringValue(""),
					"subscription_id":     types.StringValue("00000000-0000-0000-0000-000000000000"),
					"provider_namespace":  types.StringValue("Microsoft.Resources"),
					"parts": types.MapValueMust(types.StringType, map[string]attr.Value{
						"subscriptions": types.StringValue("00000000-0000-0000-0000-000000000000"),
					}),
				})),
			},
		},
		"valid-management-group-resource-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Management/managementGroups"),
					types.StringValue("/providers/Microsoft.Management/managementGroups/mg1"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.ParseResourceIdResultAttrTypes, map[string]attr.Value{
					"id":                  types.StringValue("/providers/Microsoft.Management/managementGroups/mg1"),
					"type":                types.StringValue("Microsoft.Management/managementGroups"),
					"name":                types.StringValue("mg1"),
					"parent_id":           types.StringValue("/"),
					"resource_group_name": types.StringValue(""),
					"subscription_id":     types.StringValue(""),
					"provider_namespace":  types.StringValue("Microsoft.Management"),
					"parts": types.MapValueMust(types.StringType, map[string]attr.Value{
						"managementGroups": types.StringValue("mg1"),
						"providers":        types.StringValue("Microsoft.Management"),
					}),
				})),
			},
		},
		"valid-extension-resource-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Authorization/locks"),
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/providers/Microsoft.Authorization/locks/mylock"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.ParseResourceIdResultAttrTypes, map[string]attr.Value{
					"id":                  types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/providers/Microsoft.Authorization/locks/mylock"),
					"type":                types.StringValue("Microsoft.Authorization/locks"),
					"name":                types.StringValue("mylock"),
					"parent_id":           types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"),
					"resource_group_name": types.StringValue("rg1"),
					"subscription_id":     types.StringValue("00000000-0000-0000-0000-000000000000"),
					"provider_namespace":  types.StringValue("Microsoft.Authorization"),
					"parts": types.MapValueMust(types.StringType, map[string]attr.Value{
						"subscriptions":   types.StringValue("00000000-0000-0000-0000-000000000000"),
						"resourceGroups":  types.StringValue("rg1"),
						"providers":       types.StringValue("Microsoft.Authorization"),
						"virtualNetworks": types.StringValue("vnet1"),
						"locks":           types.StringValue("mylock"),
					}),
				})),
			},
		},
		"empty-resource-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Network/virtualNetworks"),
					types.StringValue(""),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("failed to parse resource ID(resourceType: Microsoft.Network/virtualNetworks@latest, resourceId: ): `resource_id` and `type` are not matched, expect `type` to be , but got Microsoft.Network/virtualNetworks"),
				Result: function.NewResultData(types.ObjectUnknown(functions.ParseResourceIdResultAttrTypes)),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := function.RunResponse{
				Result: function.NewResultData(types.ObjectUnknown(functions.ParseResourceIdResultAttrTypes)),
			}

			parseResourceIdFunction := functions.ParseResourceIdFunction{}
			parseResourceIdFunction.Run(context.Background(), testCase.request, &got)
			fmt.Printf("got: %v\n", got)
			fmt.Printf("expected: %v\n", testCase.expected)
			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
