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
			result, lastRetryErr := clients.NewRetryOptions(tc.input)
			if result == nil && tc.expected == nil {
				if lastRetryErr != nil {
					t.Errorf("expected nil lastRetryErr for nil result, got %v", lastRetryErr)
				}
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
			result, lastRetryErr := clients.NewRetryOptionsForReadAfterCreate()
			if result == nil && tc.expected == nil {
				return
			}
			if result == nil || tc.expected == nil {
				t.Errorf("expected %v, got %v", tc.expected, result)
				return
			}
			if lastRetryErr == nil {
				t.Error("expected non-nil lastRetryErr")
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
			result, _ := clients.CombineRetryOptions(tc.input...)
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

	result, _ := clients.NewRetryOptions(retryInput)
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
	result, _ := clients.NewRetryOptionsForReadAfterCreate()
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

// Integration tests: verify ShouldRetry closures capture the last retryable error into LastRetryError

func newTestRetryValue(regexes []string) retry.RetryValue {
	regexValues := make([]attr.Value, len(regexes))
	for i, r := range regexes {
		regexValues[i] = types.StringValue(r)
	}
	return retry.NewRetryValueMust(
		map[string]attr.Type{
			"interval_seconds":     types.Int64Type,
			"max_interval_seconds": types.Int64Type,
			"multiplier":           types.Float64Type,
			"randomization_factor": types.Float64Type,
			"error_message_regex":  types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"interval_seconds":     types.Int64Value(1),
			"max_interval_seconds": types.Int64Value(10),
			"multiplier":           types.Float64Value(1.5),
			"randomization_factor": types.Float64Value(0.5),
			"error_message_regex":  types.ListValueMust(types.StringType, regexValues),
		},
	)
}

func Test_NewRetryOptions_LastRetryErrorCapture(t *testing.T) {
	retryValue := newTestRetryValue([]string{"retryable error"})

	opts, lastErr := clients.NewRetryOptions(retryValue)
	if opts == nil {
		t.Fatal("expected non-nil RetryOptions")
	}
	if lastErr == nil {
		t.Fatal("expected non-nil LastRetryError")
	}

	// Initially, no error should be captured
	initialErr := lastErr.Get()
	if initialErr != nil {
		t.Errorf("expected nil initial error, got: %v", initialErr)
	}

	// Simulate a retryable error (matching the regex)
	retryableErr := errors.New("this is a retryable error from API")
	shouldRetry := opts.ShouldRetry(nil, retryableErr)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true for matching error")
	}
	capturedErr := lastErr.Get()
	if capturedErr == nil {
		t.Fatal("expected captured error after retryable error")
	}
	if capturedErr.Error() != retryableErr.Error() {
		t.Errorf("expected captured error %q, got %q", retryableErr.Error(), capturedErr.Error())
	}

	// Simulate a non-retryable error — should not overwrite
	nonRetryableErr := errors.New("some other error")
	shouldRetry = opts.ShouldRetry(nil, nonRetryableErr)
	if shouldRetry {
		t.Error("expected ShouldRetry to return false for non-matching error")
	}
	// LastRetryError should still hold the previous retryable error
	stillCaptured := lastErr.Get()
	if stillCaptured.Error() != retryableErr.Error() {
		t.Errorf("expected captured error to remain %q, got %q", retryableErr.Error(), stillCaptured.Error())
	}
}

func Test_NewRetryOptionsForReadAfterCreate_LastRetryErrorCapture(t *testing.T) {
	opts, lastErr := clients.NewRetryOptionsForReadAfterCreate()
	if opts == nil {
		t.Fatal("expected non-nil RetryOptions")
	}
	if lastErr == nil {
		t.Fatal("expected non-nil LastRetryError")
	}

	// Simulate a 404 response — should be retryable for read-after-create
	resp404 := &http.Response{StatusCode: 404}
	err404 := errors.New("resource not found")
	shouldRetry := opts.ShouldRetry(resp404, err404)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true for 404")
	}
	captured404 := lastErr.Get()
	if captured404 == nil {
		t.Fatal("expected captured error after 404")
	}

	// Simulate a 403 response — should also be retryable for read-after-create
	resp403 := &http.Response{StatusCode: 403}
	err403 := errors.New("forbidden")
	shouldRetry = opts.ShouldRetry(resp403, err403)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true for 403")
	}
	captured403 := lastErr.Get()
	if captured403 == nil {
		t.Fatal("expected captured error after 403")
	}
	if captured403.Error() != err403.Error() {
		t.Errorf("expected captured error %q, got %q", err403.Error(), captured403.Error())
	}
}

func Test_CombineRetryOptions_LastRetryErrorCapture(t *testing.T) {
	// Create two child retry options
	retryValue := newTestRetryValue([]string{"retryable"})
	childOpts1, _ := clients.NewRetryOptions(retryValue)
	childOpts2, _ := clients.NewRetryOptionsForReadAfterCreate()

	combined, lastErr := clients.CombineRetryOptions(childOpts1, childOpts2)
	if combined == nil {
		t.Fatal("expected non-nil combined RetryOptions")
	}
	if lastErr == nil {
		t.Fatal("expected non-nil LastRetryError from CombineRetryOptions")
	}

	// Trigger via error message match (child1)
	retryableErr := errors.New("this is retryable")
	shouldRetry := combined.ShouldRetry(nil, retryableErr)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true via error message match")
	}
	capturedMsg := lastErr.Get()
	if capturedMsg == nil || capturedMsg.Error() != retryableErr.Error() {
		t.Errorf("expected combined LastRetryError to capture %q, got %v", retryableErr.Error(), capturedMsg)
	}

	// Trigger via status code match (child2: 404)
	resp404 := &http.Response{StatusCode: 404}
	err404 := errors.New("not found")
	shouldRetry = combined.ShouldRetry(resp404, err404)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true via 404 status code")
	}
	captured404 := lastErr.Get()
	if captured404 == nil || captured404.Error() != err404.Error() {
		t.Errorf("expected combined LastRetryError to capture %q, got %v", err404.Error(), captured404)
	}
}

func Test_NewRetryOptions_NullInput_ReturnsNilLastRetryError(t *testing.T) {
	opts, lastErr := clients.NewRetryOptions(retry.NewRetryValueNull())
	if opts != nil {
		t.Error("expected nil RetryOptions for null input")
	}
	if lastErr != nil {
		t.Error("expected nil LastRetryError for null input")
	}
}

func Test_CombineRetryOptions_AllNilInput_ReturnsNonNilLastRetryError(t *testing.T) {
	// When called with nil opts, CombineRetryOptions still returns a combined result
	// (because the slice is non-empty — it has nil elements).
	// Only an empty variadic call returns nil.
	combined, lastErr := clients.CombineRetryOptions(nil, nil)
	if combined == nil {
		t.Error("expected non-nil combined RetryOptions for nil opts (non-empty variadic)")
	}
	if lastErr == nil {
		t.Error("expected non-nil LastRetryError for nil opts (non-empty variadic)")
	}

	// Empty variadic call should return nil
	combined2, lastErr2 := clients.CombineRetryOptions()
	if combined2 != nil {
		t.Error("expected nil combined RetryOptions for empty variadic call")
	}
	if lastErr2 != nil {
		t.Error("expected nil LastRetryError for empty variadic call")
	}
}

func Test_NewRetryOptions_SuccessResponseNotRetried(t *testing.T) {
	// Regression test: a 200 OK response should NOT be retried even with a catch-all regex like ".*"
	retryValue := newTestRetryValue([]string{".*"})
	opts, lastErr := clients.NewRetryOptions(retryValue)
	if opts == nil {
		t.Fatal("expected non-nil RetryOptions")
	}

	// Simulate a successful 200 OK response with no error — should NOT be retried
	resp200 := &http.Response{StatusCode: 200}
	shouldRetry := opts.ShouldRetry(resp200, nil)
	if shouldRetry {
		t.Error("expected ShouldRetry to return false for 200 OK response, but got true")
	}

	// LastRetryError should not have been set
	capturedErr := lastErr.Get()
	if capturedErr != nil {
		t.Errorf("expected no captured error for 200 OK response, got: %v", capturedErr)
	}

	// But a 4xx error response with nil err should still be retried via regex
	resp500 := &http.Response{StatusCode: 500}
	shouldRetry = opts.ShouldRetry(resp500, nil)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true for 500 response (matches status code)")
	}

	// A 400 error response (not in default status codes) with nil err should be retried via regex on error message
	resp400 := &http.Response{StatusCode: 400}
	shouldRetry = opts.ShouldRetry(resp400, nil)
	if !shouldRetry {
		t.Error("expected ShouldRetry to return true for 400 response with .* regex")
	}
}
