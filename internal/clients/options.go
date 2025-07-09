package clients

import (
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
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

	statusCodes := make([]int, 0)
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		statusCodes = append(statusCodes, opt.StatusCodes...)
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
			if opt == nil {
				continue
			}
			if opt.ShouldRetry != nil {
				if opt.ShouldRetry(resp, err) {
					return true
				}
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
func NewRetryOptionsForReadAfterCreate(maxRetries int32) *policy.RetryOptions {
	statusCodes := make([]int, 0)
	statusCodes = append(statusCodes, DefaultRetryableStatusCodes...)
	// Add default read after create values for the default retry configuration.
	statusCodes = append(statusCodes, DefaultRetryableReadAfterCreateStatusCodes...)
	return &policy.RetryOptions{
		MaxRetries:  maxRetries,
		StatusCodes: statusCodes,
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
			if err != nil {
				// Check if the error message matches any of the retryable error regexps.
				for _, re := range rtry.GetErrorMessagesRegex() {
					if re.MatchString(err.Error()) {
						return true
					}
				}
			} else {
				responseErr := runtime.NewResponseError(resp)
				if responseErr != nil {
					// Check if the error message matches any of the retryable error regexps.
					for _, re := range rtry.GetErrorMessagesRegex() {
						if re.MatchString(responseErr.Error()) {
							log.Printf("[DEBUG] Retrying request due to error: %s matches regex %s", responseErr.Error(), re.String())
							return true
						}
					}
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
