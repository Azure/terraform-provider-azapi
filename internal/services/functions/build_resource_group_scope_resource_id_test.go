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

func Test_ResourceGroupResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"resource-group-scope-valid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("rg1"),
					types.StringValue("Microsoft.Network/virtualNetworks"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("vnet1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1")),
			},
		},
		"resource-group-scope-valid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("rg1"),
					types.StringValue("Microsoft.Network/virtualNetworks/subnets"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("vnet1"),
						types.StringValue("subnet1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1")),
			},
		},
		"resource-group-scope-invalid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("rg1"),
					types.StringValue("Microsoft.Network/virtualNetworks"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("vnet1"),
						types.StringValue("subnet1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 1, got 2"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"resource-group-scope-invalid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("rg1"),
					types.StringValue("Microsoft.Network/virtualNetworks/subnets"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("vnet1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 2, got 1"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"resource-group-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("rg1"),
					types.StringValue(""),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("vnet1"),
						types.StringValue("subnet1"),
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

			resourceGroupResourceIdFunction := functions.ResourceGroupResourceIdFunction{}
			resourceGroupResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
