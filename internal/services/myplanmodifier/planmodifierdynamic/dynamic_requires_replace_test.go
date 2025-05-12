package planmodifierdynamic

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RequiresReplaceIfNotNull(t *testing.T) {
	tests := []struct {
		name           string
		planValue      types.Dynamic
		stateValue     types.Dynamic
		expectedResult bool
	}{
		{
			name:           "both values null",
			planValue:      types.DynamicNull(),
			stateValue:     types.DynamicNull(),
			expectedResult: false,
		},
		{
			name:           "plan value null",
			planValue:      types.DynamicNull(),
			stateValue:     types.DynamicValue(types.StringValue("state")),
			expectedResult: false,
		},
		{
			name:           "state value null",
			planValue:      types.DynamicValue(types.StringValue("plan")),
			stateValue:     types.DynamicNull(),
			expectedResult: false,
		},
		{
			name:           "values equal",
			planValue:      types.DynamicValue(types.StringValue("same")),
			stateValue:     types.DynamicValue(types.StringValue("same")),
			expectedResult: false,
		},
		{
			name:           "values different",
			planValue:      types.DynamicValue(types.StringValue("plan")),
			stateValue:     types.DynamicValue(types.StringValue("state")),
			expectedResult: true,
		},
	}

	// We need to create a dummy state and plan because the PlanModifierDynamic
	// method checks for a null state and plan, returning early if they are null.
	state := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"test": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"test": tftypes.NewValue(tftypes.String, "state"),
	})

	plan := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"test": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"test": tftypes.NewValue(tftypes.String, "plan"),
	})

	// Set up the plan and state values for the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := RequiresReplaceIfNotNull()

			req := planmodifier.DynamicRequest{
				State: tfsdk.State{
					Raw: state,
				},
				Plan: tfsdk.Plan{
					Raw: plan,
				},
				PlanValue:  tt.planValue,
				StateValue: tt.stateValue,
			}

			resp := &planmodifier.DynamicResponse{}

			modifier.PlanModifyDynamic(context.Background(), req, resp)

			require.Empty(t, resp.Diagnostics)
			assert.Equal(t, tt.expectedResult, resp.RequiresReplace)
		})
	}
}
