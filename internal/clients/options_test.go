package clients_test

import (
	"errors"
	"math"
	"net/http"
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

func Test_NewRetryOptions_NilResponseHandling(t *testing.T) {
	// This test reproduces the panic from issue #985
	// where resp.StatusCode is accessed when resp is nil
	retryInput := retry.NewRetryValueMust(
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
	)

	result := clients.NewRetryOptions(retryInput)
	if result == nil {
		t.Fatal("expected non-nil retry options")
	}

	// Test various scenarios
	testCases := []struct {
		name        string
		resp        *http.Response
		err         error
		expectRetry bool
	}{
		{
			name:        "nil response with matching error",
			resp:        nil,
			err:         errors.New("some network error"),
			expectRetry: true, // Should retry because error message matches regex
		},
		{
			name:        "nil response with non-matching error",
			resp:        nil,
			err:         errors.New("some other issue"),
			expectRetry: false,
		},
		{
			name:        "retryable status code",
			resp:        &http.Response{StatusCode: 429}, // Too Many Requests
			err:         nil,
			expectRetry: true,
		},
		{
			name:        "non-retryable status code",
			resp:        &http.Response{StatusCode: 400}, // Bad Request
			err:         nil,
			expectRetry: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShouldRetry panicked: %v", r)
				}
			}()

			shouldRetry := result.ShouldRetry(tc.resp, tc.err)
			if shouldRetry != tc.expectRetry {
				t.Errorf("expected ShouldRetry to return %v, got %v", tc.expectRetry, shouldRetry)
			}
		})
	}
}

func Test_NewRetryOptionsForReadAfterCreate_NilResponseHandling(t *testing.T) {
	// This test reproduces the panic from issue #985 for the read-after-create retry options
	result := clients.NewRetryOptionsForReadAfterCreate()
	if result == nil {
		t.Fatal("expected non-nil retry options")
	}

	// Test various scenarios
	testCases := []struct {
		name        string
		resp        *http.Response
		err         error
		expectRetry bool
	}{
		{
			name:        "nil response with error",
			resp:        nil,
			err:         errors.New("some network error"),
			expectRetry: false, // This function doesn't check error messages, only status codes
		},
		{
			name:        "retryable status code",
			resp:        &http.Response{StatusCode: 429}, // Too Many Requests
			err:         nil,
			expectRetry: true,
		},
		{
			name:        "non-retryable status code",
			resp:        &http.Response{StatusCode: 400}, // Bad Request
			err:         nil,
			expectRetry: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShouldRetry panicked: %v", r)
				}
			}()

			shouldRetry := result.ShouldRetry(tc.resp, tc.err)
			if shouldRetry != tc.expectRetry {
				t.Errorf("expected ShouldRetry to return %v, got %v", tc.expectRetry, shouldRetry)
			}
		})
	}
}
