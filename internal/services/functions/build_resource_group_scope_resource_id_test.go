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

func TestBuildResourceGroupScopeResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"resource-group-scope-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("myResourceGroup"),
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
		"resource-group-scope-invalid-empty-subscription-id": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(""),
					types.StringValue("myResourceGroup"),
					types.StringValue("Microsoft.Automation/automationAccounts@2021-06-22"),
					types.StringValue("myAutomationAccount"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("subscription_id cannot be empty"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"resource-group-scope-invalid-empty-resource-group-name": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue(""),
					types.StringValue("Microsoft.Automation/automationAccounts@2021-06-22"),
					types.StringValue("myAutomationAccount"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("resource_group_name cannot be empty"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"resource-group-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("myResourceGroup"),
					types.StringValue(""),
					types.StringValue("myAutomationAccount"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`type` is invalid, it should be like `ResourceProvider/resourceTypes@ApiVersion`"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"resource-group-scope-invalid-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("00000000-0000-0000-0000-000000000000"),
					types.StringValue("myResourceGroup"),
					types.StringValue("Invalid/ResourceType"),
					types.StringValue("myAutomationAccount"),
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

			resourceGroupScopeResourceIdFunction := functions.BuildResourceGroupScopeResourceIdFunction{}
			resourceGroupScopeResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
