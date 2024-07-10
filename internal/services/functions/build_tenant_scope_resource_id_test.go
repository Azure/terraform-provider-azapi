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

func TestBuildTenantScopeResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"tenant-scope-valid": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
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
		"tenant-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(""),
					types.StringValue("MyInvoice"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("`type` is invalid, it should be like `ResourceProvider/resourceTypes@ApiVersion`"),
				Result: function.NewResultData(types.ObjectUnknown(functions.BuildResourceIdResultAttrTypes)),
			},
		},
		"tenant-scope-invalid-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Invalid/ResourceType"),
					types.StringValue("MyInvoice"),
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

			tenantScopeResourceIdFunction := functions.BuildTenantScopeResourceIdFunction{}
			tenantScopeResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
