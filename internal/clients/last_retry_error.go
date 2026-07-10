package clients

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type LastRetryError struct {
	mu  sync.Mutex
	err error
}

func (e *LastRetryError) Set(resp *http.Response, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if err != nil {
		e.err = err
	} else if resp != nil {
		e.err = runtime.NewResponseError(resp)
	}
}

func (e *LastRetryError) Get() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.err
}

func WrapContextError(err error, lastRetryErr *LastRetryError) error {
	if err == nil {
		return nil
	}
	if lastRetryErr == nil {
		return err
	}

	if err != context.DeadlineExceeded && err != context.Canceled {
		return err
	}

	lastErr := lastRetryErr.Get()

	if lastErr != nil {
		return fmt.Errorf("%w, last retryable error: %v", err, lastErr)
	}

	return err
}
