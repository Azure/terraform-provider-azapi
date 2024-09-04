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

func Test_ExtensionResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"extension-scope-valid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM"),
					types.StringValue("Microsoft.Chaos/targets"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("t1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM/providers/Microsoft.Chaos/targets/t1")),
			},
		},
		"extension-scope-valid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM"),
					types.StringValue("Microsoft.Chaos/targets/capabilities"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("t1"),
						types.StringValue("c1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM/providers/Microsoft.Chaos/targets/t1/capabilities/c1")),
			},
		},
		"extension-scope-invalid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM"),
					types.StringValue("Microsoft.Chaos/targets"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("t1"),
						types.StringValue("c1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 1, got 2"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"extension-scope-invalid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM"),
					types.StringValue("Microsoft.Chaos/targets/capabilities"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("t1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 2, got 1"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"extension-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"),
					types.StringValue(""),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("mylock"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("invalid azure resource type, it should be like `ResourceProvider/resourceTypes`"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := function.RunResponse{
				Result: function.NewResultData(types.StringUnknown()),
			}

			extensionResourceIdFunction := functions.ExtensionResourceIdFunction{}
			extensionResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
