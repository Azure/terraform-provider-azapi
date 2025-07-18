package clients_test

import (
	"math"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_NewRetryOptions(t *testing.T) {
	testcases := []struct {
		name     string
		input    retry.RetryValue
		expected *policy.RetryOptions
	}{
		{
			name:     "null input",
			input:    retry.NewRetryValueNull(),
			expected: nil,
		},

		{
			name:     "unknown input",
			input:    retry.NewRetryValueUnknown(),
			expected: nil,
		},

		{
			name: "basic retry options",
			input: retry.NewRetryValueMust(
				map[string]attr.Type{
					"interval_seconds":     types.Int64Type,
					"max_interval_seconds": types.Int64Type,
					"multiplier":           types.Float64Type,
					"randomization_factor": types.Float64Type,
					"error_message_regex":  types.ListType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"interval_seconds":     types.Int64Value(retry.DefaultIntervalSeconds),
					"max_interval_seconds": types.Int64Value(retry.DefaultMaxIntervalSeconds),
					"multiplier":           types.Float64Value(retry.DefaultMultiplier),
					"randomization_factor": types.Float64Value(retry.DefaultRandomizationFactor),
					"error_message_regex":  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("error")}),
				},
			),
			expected: &policy.RetryOptions{
				MaxRetries:    math.MaxInt16,
				RetryDelay:    time.Second * retry.DefaultIntervalSeconds,
				MaxRetryDelay: time.Second * retry.DefaultMaxIntervalSeconds,
				StatusCodes:   clients.DefaultRetryableStatusCodes,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := clients.NewRetryOptions(tc.input)
			if result == nil && tc.expected == nil {
				return
			}
			if result == nil || tc.expected == nil {
				t.Errorf("expected %v, got %v", tc.expected, result)
				return
			}
			if result.MaxRetries != tc.expected.MaxRetries {
				t.Errorf("expected MaxRetries %d, got %d", tc.expected.MaxRetries, result.MaxRetries)
			}
			if result.RetryDelay != tc.expected.RetryDelay {
				t.Errorf("expected RetryDelay %v, got %v", tc.expected.RetryDelay, result.RetryDelay)
			}
			if result.MaxRetryDelay != tc.expected.MaxRetryDelay {
				t.Errorf("expected MaxRetryDelay %v, got %v", tc.expected.MaxRetryDelay, result.MaxRetryDelay)
			}
			if len(result.StatusCodes) != len(tc.expected.StatusCodes) {
				t.Errorf("expected %d status codes, got %d", len(tc.expected.StatusCodes), len(result.StatusCodes))
				return
			}
			for i, code := range result.StatusCodes {
				if code != tc.expected.StatusCodes[i] {
					t.Errorf("expected status code %d at index %d, got %d", tc.expected.StatusCodes[i], i, code)
				}
			}
		})
	}
}

func Test_NewRetryOptionsForReadAfterCreate(t *testing.T) {
	testcases := []struct {
		name     string
		input    int32
		expected *policy.RetryOptions
	}{
		{
			name: "default read after create",
			expected: &policy.RetryOptions{
				MaxRetries:  math.MaxInt16,
				StatusCodes: append(clients.DefaultRetryableStatusCodes, clients.DefaultRetryableReadAfterCreateStatusCodes...),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := clients.NewRetryOptionsForReadAfterCreate()
			if result == nil && tc.expected == nil {
				return
			}
			if result == nil || tc.expected == nil {
				t.Errorf("expected %v, got %v", tc.expected, result)
				return
			}
			if result.MaxRetries != tc.expected.MaxRetries {
				t.Errorf("expected MaxRetries %d, got %d", tc.expected.MaxRetries, result.MaxRetries)
			}
			if len(result.StatusCodes) != len(tc.expected.StatusCodes) {
				t.Errorf("expected %d status codes, got %d", len(tc.expected.StatusCodes), len(result.StatusCodes))
				return
			}
			for i, code := range result.StatusCodes {
				if code != tc.expected.StatusCodes[i] {
					t.Errorf("expected status code %d at index %d, got %d", tc.expected.StatusCodes[i], i, code)
				}
			}
		})
	}
}

func Test_CombineRetryOptions(t *testing.T) {
	testcases := []struct {
		name     string
		input    []*policy.RetryOptions
		expected *policy.RetryOptions
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty input",
			input:    []*policy.RetryOptions{},
			expected: nil,
		},
		{
			name: "single option",
			input: []*policy.RetryOptions{
				{
					MaxRetries:    3,
					RetryDelay:    time.Second * 2,
					MaxRetryDelay: time.Second * 5,
					StatusCodes:   []int{500, 502},
				},
			},
			expected: &policy.RetryOptions{
				MaxRetries:    3,
				RetryDelay:    time.Second * 2,
				MaxRetryDelay: time.Second * 5,
				StatusCodes:   []int{500, 502},
			},
		},
		{
			name: "multiple options",
			input: []*policy.RetryOptions{
				{
					MaxRetries:    2,
					RetryDelay:    time.Second * 2,
					MaxRetryDelay: time.Second * 10,
					StatusCodes:   []int{500, 502},
				},
				{
					MaxRetries:    5,
					RetryDelay:    time.Second * 1,
					MaxRetryDelay: time.Second * 5,
					StatusCodes:   []int{503, 504},
				},
			},
			expected: &policy.RetryOptions{
				MaxRetries:    5,                         // max of 2 and 5
				RetryDelay:    time.Second * 1,           // min of 2 and 1
				MaxRetryDelay: time.Second * 10,          // max of 10 and 5
				StatusCodes:   []int{500, 502, 503, 504}, // combined status codes
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := clients.CombineRetryOptions(tc.input...)
			if result == nil && tc.expected == nil {
				return
			}
			if result == nil || tc.expected == nil {
				t.Errorf("expected %v, got %v", tc.expected, result)
				return
			}
			if result.MaxRetries != tc.expected.MaxRetries {
				t.Errorf("expected MaxRetries %d, got %d", tc.expected.MaxRetries, result.MaxRetries)
			}
			if result.RetryDelay != tc.expected.RetryDelay {
				t.Errorf("expected RetryDelay %v, got %v", tc.expected.RetryDelay, result.RetryDelay)
			}
			if result.MaxRetryDelay != tc.expected.MaxRetryDelay {
				t.Errorf("expected MaxRetryDelay %v, got %v", tc.expected.MaxRetryDelay, result.MaxRetryDelay)
			}
			if len(result.StatusCodes) != len(tc.expected.StatusCodes) {
				t.Errorf("expected %d status codes, got %d", len(tc.expected.StatusCodes), len(result.StatusCodes))
				return
			}
			for i, code := range result.StatusCodes {
				found := false
				for _, expectedCode := range tc.expected.StatusCodes {
					if code == expectedCode {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected status code %d at index %d, but it was not found in the expected codes", code, i)
					continue
				}
			}
		})
	}
}
