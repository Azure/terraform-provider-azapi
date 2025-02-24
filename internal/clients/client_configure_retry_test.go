package clients

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestConfigureCustomRetry(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name                    string
		rtry                    retry.RetryValue
		useReadAfterCreate      bool
		expectedInitialInterval time.Duration
		expectedMaxInterval     time.Duration
		expectedMultiplier      float64
		expectedRandomization   float64
		expectedStatusCodes     []int
		expectedErrorRegexps    []string
		expectedDataCallbacks   int
	}{
		{
			name:                    "default retry configuration",
			rtry:                    retry.NewRetryValueNull(),
			useReadAfterCreate:      false,
			expectedInitialInterval: 5 * time.Second,
			expectedMaxInterval:     30 * time.Second,
			expectedMultiplier:      1.5,
			expectedRandomization:   0.5,
			expectedStatusCodes:     []int{429},
			expectedErrorRegexps:    []string{},
			expectedDataCallbacks:   0,
		},
		{
			name:                    "default retry configuration with read after create",
			rtry:                    retry.NewRetryValueNull(),
			useReadAfterCreate:      true,
			expectedInitialInterval: 5 * time.Second,
			expectedMaxInterval:     30 * time.Second,
			expectedMultiplier:      1.5,
			expectedRandomization:   0.5,
			expectedStatusCodes:     []int{404, 403},
			expectedErrorRegexps:    []string{},
			expectedDataCallbacks:   0,
		},
		{
			name: "custom retry configuration",
			rtry: retry.NewRetryValueMust(retry.NewRetryValueNull().AttributeTypes(ctx), map[string]attr.Value{
				"interval_seconds":     basetypes.NewInt64Value(10),
				"max_interval_seconds": basetypes.NewInt64Value(60),
				"multiplier":           basetypes.NewFloat64Value(1.5),
				"randomization_factor": basetypes.NewFloat64Value(0.2),
				"http_status_codes":    basetypes.NewListValueMust(types.Int64Type, []attr.Value{basetypes.NewInt64Value(500), basetypes.NewInt64Value(503)}),
				"error_message_regex":  basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("timeout"), basetypes.NewStringValue("temporary")}),
				"read_after_create_retry_on_not_found_or_forbidden": basetypes.NewBoolValue(true),
				"response_is_nil": basetypes.NewBoolValue(false),
			}),
			useReadAfterCreate:      false,
			expectedInitialInterval: 10 * time.Second,
			expectedMaxInterval:     60 * time.Second,
			expectedMultiplier:      1.5,
			expectedRandomization:   0.2,
			expectedStatusCodes:     []int{500, 503},
			expectedErrorRegexps:    []string{"timeout", "temporary"},
			expectedDataCallbacks:   0,
		},
		{
			name: "custom retry with read after create",
			rtry: retry.NewRetryValueMust(retry.NewRetryValueNull().AttributeTypes(ctx), map[string]attr.Value{
				"interval_seconds":     basetypes.NewInt64Value(10),
				"max_interval_seconds": basetypes.NewInt64Value(60),
				"multiplier":           basetypes.NewFloat64Value(1.5),
				"randomization_factor": basetypes.NewFloat64Value(0.2),
				"http_status_codes":    basetypes.NewListValueMust(types.Int64Type, []attr.Value{basetypes.NewInt64Value(500), basetypes.NewInt64Value(503)}),
				"error_message_regex":  basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("timeout"), basetypes.NewStringValue("temporary")}),
				"read_after_create_retry_on_not_found_or_forbidden": basetypes.NewBoolValue(true),
				"response_is_nil": basetypes.NewBoolValue(false),
			}),
			useReadAfterCreate:      true,
			expectedInitialInterval: 10 * time.Second,
			expectedMaxInterval:     60 * time.Second,
			expectedMultiplier:      1.5,
			expectedRandomization:   0.2,
			expectedStatusCodes:     []int{500, 503, 404, 403},
			expectedErrorRegexps:    []string{"timeout", "temporary"},
			expectedDataCallbacks:   0,
		},
		{
			name: "custom retry with response is nil",
			rtry: retry.NewRetryValueMust(retry.NewRetryValueNull().AttributeTypes(ctx), map[string]attr.Value{
				"interval_seconds":     basetypes.NewInt64Value(10),
				"max_interval_seconds": basetypes.NewInt64Value(60),
				"multiplier":           basetypes.NewFloat64Value(1.5),
				"randomization_factor": basetypes.NewFloat64Value(0.2),
				"http_status_codes":    basetypes.NewListValueMust(types.Int64Type, []attr.Value{basetypes.NewInt64Value(500), basetypes.NewInt64Value(503)}),
				"error_message_regex":  basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("timeout"), basetypes.NewStringValue("temporary")}),
				"read_after_create_retry_on_not_found_or_forbidden": basetypes.NewBoolValue(true),
				"response_is_nil": basetypes.NewBoolValue(true),
			}),
			useReadAfterCreate:      false,
			expectedInitialInterval: 10 * time.Second,
			expectedMaxInterval:     60 * time.Second,
			expectedMultiplier:      1.5,
			expectedRandomization:   0.2,
			expectedStatusCodes:     []int{500, 503},
			expectedErrorRegexps:    []string{"timeout", "temporary"},
			expectedDataCallbacks:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			backOff, errRegExps, statusCodes, dataCallbackFuncs := configureCustomRetry(ctx, tt.rtry, tt.useReadAfterCreate)

			assert.Equalf(t, tt.expectedInitialInterval, backOff.InitialInterval, "InitialInterval")
			assert.Equalf(t, tt.expectedMaxInterval, backOff.MaxInterval, "MaxInterval")
			assert.Equalf(t, tt.expectedMultiplier, backOff.Multiplier, "Multiplier")
			assert.Equalf(t, tt.expectedRandomization, backOff.RandomizationFactor, "RandomizationFactor")
			assert.Equalf(t, tt.expectedStatusCodes, statusCodes, "StatusCodes")

			actualErrorRegexps := make([]string, len(errRegExps))
			for i, re := range errRegExps {
				actualErrorRegexps[i] = re.String()
			}
			assert.Equal(t, tt.expectedErrorRegexps, actualErrorRegexps, "ErrorRegexps")
			assert.Len(t, dataCallbackFuncs, tt.expectedDataCallbacks, "DataCallbacks")
		})
	}
}
