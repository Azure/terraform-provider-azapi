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

func Test_TenantResourceIdFunction(t *testing.T) {
	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"tenant-scope-valid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Billing/billingAccounts"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("ba1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Billing/billingAccounts/ba1")),
			},
		},
		"tenant-scope-valid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Billing/billingAccounts/billingProfiles"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("ba1"),
						types.StringValue("bp1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Billing/billingAccounts/ba1/billingProfiles/bp1")),
			},
		},
		"tenant-scope-valid-triple": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Billing/billingAccounts/billingProfiles/invoiceSections@2020-05-01"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("ba1"),
						types.StringValue("bp1"),
						types.StringValue("is1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("/providers/Microsoft.Billing/billingAccounts/ba1/billingProfiles/bp1/invoiceSections/is1")),
			},
		},
		"tenant-scope-invalid-single": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Billing/billingAccounts"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("ba1"),
						types.StringValue("bp1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 1, got 2"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"tenant-scope-invalid-double": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("Microsoft.Billing/billingAccounts/billingProfiles"),
					types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("ba1"),
					}),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("number of resource names does not match the number of resource type parts, expected 2, got 1"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"tenant-scope-invalid-empty-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
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

			tenantResourceIdFunction := functions.TenantResourceIdFunction{}
			tenantResourceIdFunction.Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
