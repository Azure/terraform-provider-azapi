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

func Test_ManagementGroupResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"management-group-scope-valid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue("Microsoft.Blueprint/blueprints"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("bp1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Blueprint/blueprints/bp1")),
			},
		},
		"management-group-scope-valid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue("Microsoft.Blueprint/blueprints/artifacts"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("bp1"),
						types.StringValue("a1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Blueprint/blueprints/bp1/artifacts/a1")),
			},
		},
		"management-group-scope-invalid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue("Microsoft.Blueprint/blueprints"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("bp1"),
						types.StringValue("a1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 1, got 2"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"management-group-scope-invalid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue("Microsoft.Blueprint/blueprints/artifacts"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("bp1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 2, got 1"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"management-group-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue(""),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("ba1"),
						types.StringValue("bp1"),
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

			managementGroupResourceIdFunction := functions.ManagementGroupResourceIdFunction{}
			managementGroupResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
