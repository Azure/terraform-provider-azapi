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

func TestBuildManagementGroupScopeResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"management-group-scope-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue("Microsoft.CostManagement/views@2023-04-01-preview"),
					types.StringValue("MyCostView"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.ObjectValueMust(functions.BuildResourceIdResultAttrTypes, map[string]attr.Value{
					"resource_id": types.StringValue("/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.CostManagement/views/MyCostView"),
				})),
			},
		},
		"management-group-scope-invalid-empty-mg-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(""),
					types.StringValue("Microsoft.CostManagement/views@2023-04-01-preview"),
					types.StringValue("MyCostView"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("management_group_id cannot be empty"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"management-group-scope-invalid-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue("Invalid/ResourceType"),
					types.StringValue("name"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`type` is invalid, it should be like `ResourceProvider/resourceTypes@ApiVersion`"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"management-group-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("mg1"),
					types.StringValue(""),
					types.StringValue("MyCostView"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`type` is invalid, it should be like `ResourceProvider/resourceTypes@ApiVersion`"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := function.RunResponse{
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			}

			managementGroupScopeResourceIdFunction := functions.BuildManagementGroupScopeResourceIdFunction{}
			managementGroupScopeResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
