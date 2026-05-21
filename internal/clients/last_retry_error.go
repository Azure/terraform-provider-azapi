package clients

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// LastRetryError captures the last retryable error encountered during retry loops.
// It is populated by the ShouldRetry closure and read after the pipeline call returns.
// This allows us to surface the actual API error instead of the unhelpful
// "context deadline exceeded" message when the retry loop exhausts the context timeout.
type LastRetryError struct {
	mu  sync.Mutex
	err error
}

// Set records the last retryable error. It is called from ShouldRetry closures.
// When err is nil but resp is non-nil, a ResponseError is constructed from the
// response to capture the full ARM error body (status code, error code, message).
func (e *LastRetryError) Set(resp *http.Response, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if err != nil {
		e.err = err
	} else if resp != nil {
		// Construct a ResponseError from the response to capture the full ARM error text.
		// The response body is safely cached by the SDK's bodyDownloadPolicy, so
		// NewResponseError can read and parse the JSON/XML error body.
		e.err = runtime.NewResponseError(resp)
	}
}

// Get returns the last captured retryable error.
func (e *LastRetryError) Get() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.err
}

// WrapContextError checks if err is a context deadline or cancellation error,
// and if so, wraps it with the last retryable error captured by LastRetryError.
// If lastRetryErr is nil or has no captured error, the original error is returned unchanged.
func WrapContextError(err error, lastRetryErr *LastRetryError) error {
	if err == nil {
		return nil
	}
	if lastRetryErr == nil {
		return err
	}

	// Only wrap context-related errors
	if err != context.DeadlineExceeded && err != context.Canceled {
		return err
	}

	lastErr := lastRetryErr.Get()

	// If we have a captured error, wrap the context error with it
	if lastErr != nil {
		return fmt.Errorf("%w, last retryable error: %v", err, lastErr)
	}

	return err
}
