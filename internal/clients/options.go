package clients

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	jmes "github.com/jmespath/go-jmespath"
)

var DefaultRetryableStatusCodes = []int{
	http.StatusRequestTimeout,      // 408
	http.StatusTooManyRequests,     // 429
	http.StatusInternalServerError, // 500
	http.StatusBadGateway,          // 502
	http.StatusServiceUnavailable,  // 503
	http.StatusGatewayTimeout,      // 504
}
var DefaultRetryableReadAfterCreateStatusCodes = []int{
	http.StatusNotFound,  // 404
	http.StatusForbidden, // 403
}

type RequestOptions struct {
	Headers         map[string]string
	QueryParameters map[string]string
	RetryOptions    *policy.RetryOptions
}

// CombineRetryOptions combines multiple RequestOptions into a single policy.RetryOptions.
func CombineRetryOptions(opts ...*policy.RetryOptions) *policy.RetryOptions {
	if len(opts) == 0 {
		return nil
	}

	statusCodeSet := make(map[int]bool)
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		for _, code := range opt.StatusCodes {
			statusCodeSet[code] = true
		}
	}

	statusCodes := make([]int, 0)
	for code := range statusCodeSet {
		statusCodes = append(statusCodes, code)
	}

	var maxRetries int32 = 0
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.MaxRetries > maxRetries {
			maxRetries = opt.MaxRetries
		}
	}

	var retryDelay time.Duration = 0
	var maxRetryDelay time.Duration = 0
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if retryDelay == 0 || opt.RetryDelay < retryDelay {
			retryDelay = opt.RetryDelay
		}
		if opt.MaxRetryDelay > maxRetryDelay {
			maxRetryDelay = opt.MaxRetryDelay
		}
	}

	shouldRetry := func(resp *http.Response, err error) bool {
		for _, opt := range opts {
			if opt == nil || opt.ShouldRetry == nil {
				continue
			}
			if opt.ShouldRetry(resp, err) {
				return true
			}
		}
		return false
	}

	return &policy.RetryOptions{
		MaxRetries:    maxRetries,
		RetryDelay:    retryDelay,
		MaxRetryDelay: maxRetryDelay,
		StatusCodes:   statusCodes,
		ShouldRetry:   shouldRetry,
	}
}

// NewRetryOptionsForReadAfterCreate creates a RetryOptions for read-after-create operations.
func NewRetryOptionsForReadAfterCreate() *policy.RetryOptions {
	log.Printf("[DEBUG] Using custom retry configuration for read after create")
	statusCodes := make([]int, 0)
	statusCodes = append(statusCodes, DefaultRetryableStatusCodes...)
	// Add default read after create values for the default retry configuration.
	statusCodes = append(statusCodes, DefaultRetryableReadAfterCreateStatusCodes...)
	return &policy.RetryOptions{
		// Set a very high max retries to make sure context deadline is respected.
		MaxRetries:  math.MaxInt16,
		StatusCodes: statusCodes,
		ShouldRetry: func(resp *http.Response, err error) bool {
			// We need to test for status codes here too. This covers the case that these options are combined with
			// retry options from NewRetryOptions, because the ShouldRetry function takes precedence over StatusCodes.
			if resp != nil {
				for _, code := range statusCodes {
					if resp.StatusCode == code {
						return true
					}
				}
			}
			return false
		},
	}
}

// NewRetryOptionsForWaitForDesiredState creates a RetryOptions that evaluates JMESPath expressions
// against the response body after a successful GET request. If any expression does not evaluate to true,
// the request will be retried.
func NewRetryOptionsForWaitForDesiredState(expressions []string) *policy.RetryOptions {
	if len(expressions) == 0 {
		return nil
	}

	log.Printf("[DEBUG] Using custom retry configuration for wait for desired state with %d expression(s)", len(expressions))
	return &policy.RetryOptions{
		// Set a very high max retries to make sure context deadline is respected.
		MaxRetries:  math.MaxInt16,
		StatusCodes: DefaultRetryableStatusCodes,
		ShouldRetry: func(resp *http.Response, err error) bool {
			if resp == nil || resp.Body == nil {
				return false
			}

			// Only evaluate JMESPath expressions for successful responses
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return false
			}

			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				log.Printf("[DEBUG] Failed to read response body for desired state check: %s", readErr.Error())
				// Replace the body so it can be read by the caller
				resp.Body = io.NopCloser(bytes.NewReader(nil))
				return false
			}

			// Always replace the body so it can be read by the caller
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			var body interface{}
			if jsonErr := json.Unmarshal(bodyBytes, &body); jsonErr != nil {
				log.Printf("[DEBUG] Failed to parse response body as JSON for desired state check: %s", jsonErr.Error())
				return false
			}

			for _, expr := range expressions {
				result, searchErr := jmes.Search(expr, body)
				if searchErr != nil {
					log.Printf("[WARN] Failed to evaluate JMESPath expression %q for desired state check: %s", expr, searchErr.Error())
					return false
				}

				boolResult, ok := result.(bool)
				if !ok || !boolResult {
					log.Printf("[DEBUG] Retrying: JMESPath expression %q did not evaluate to true (got: %v)", expr, result)
					return true
				}
			}

			log.Printf("[DEBUG] All JMESPath expressions evaluated to true, desired state reached")
			return false
		},
	}
}

// NewRetryOptions creates a RetryOptions based on the provided retry.RetryValue.
func NewRetryOptions(rtry retry.RetryValue) *policy.RetryOptions {
	if rtry.IsNull() || rtry.IsUnknown() {
		return nil
	}

	log.Printf("[DEBUG] Using custom retry configuration")
	return &policy.RetryOptions{
		RetryDelay:    rtry.GetIntervalSecondsAsDuration(),
		MaxRetryDelay: rtry.GetMaxIntervalSecondsAsDuration(),
		// Set a very high max retries to make sure context deadline is respected.
		MaxRetries:  math.MaxInt16,
		StatusCodes: DefaultRetryableStatusCodes,
		ShouldRetry: func(resp *http.Response, err error) bool {
			// We need to test for DefaultRetryableStatusCodes here as using ShouldRetry overrides the use of StatusCodes.
			if resp != nil {
				for _, code := range DefaultRetryableStatusCodes {
					if resp.StatusCode == code {
						return true
					}
				}
			}

			// Get the error message to check against regex patterns,
			// If use the err.Error() string first, else get the response error from the HTTP response.
			var errorMsg string
			if err != nil {
				errorMsg = err.Error()
			} else if resp != nil {
				responseErr := runtime.NewResponseError(resp)
				if responseErr != nil {
					errorMsg = responseErr.Error()
				}
			}
			// Check if the error message matches any of the retryable error regexps
			if errorMsg == "" {
				return false
			}
			for _, re := range rtry.GetErrorMessagesRegex() {
				if re.MatchString(errorMsg) {
					log.Printf("[DEBUG] Retrying request due to error: %s matches regex %s", errorMsg, re.String())
					return true
				}
			}
			return false
		},
	}
}

func NewQueryParameters(queryParameters map[string][]string) map[string]string {
	opts := make(map[string]string)

	for key, values := range queryParameters {
		if len(values) > 0 {
			opts[key] = strings.Join(values, ",")
		}
	}

	return opts
}

func DefaultRequestOptions() RequestOptions {
	return RequestOptions{
		Headers:         make(map[string]string),
		QueryParameters: make(map[string]string),
	}
}

func NewRequestOptions(headers map[string]string, queryParameters map[string][]string) RequestOptions {
	opts := DefaultRequestOptions()

	opts.Headers = headers

	for key, values := range queryParameters {
		opts.QueryParameters[key] = strings.Join(values, ",")
	}

	return opts
}
