package clients_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
)

func Test_LastRetryError_SetAndGet(t *testing.T) {
	lre := &clients.LastRetryError{}

	// Initially should return nil
	err := lre.Get()
	if err != nil {
		t.Fatalf("expected nil; got %v", err)
	}

	// Set with an explicit error
	testErr := errors.New("test error")
	testResp := &http.Response{StatusCode: http.StatusInternalServerError}
	lre.Set(testResp, testErr)

	err = lre.Get()
	// When err is non-nil, Set stores it directly
	if err != testErr {
		t.Errorf("expected error %v, got %v", testErr, err)
	}
}

func Test_LastRetryError_OverwritesPrevious(t *testing.T) {
	lre := &clients.LastRetryError{}

	firstErr := errors.New("first error")
	lre.Set(nil, firstErr)

	secondErr := errors.New("second error")
	secondResp := &http.Response{StatusCode: http.StatusTooManyRequests}
	lre.Set(secondResp, secondErr)

	err := lre.Get()
	if err != secondErr {
		t.Errorf("expected error %v, got %v", secondErr, err)
	}
}

func Test_LastRetryError_SetWithNilError_ConstructsResponseError(t *testing.T) {
	lre := &clients.LastRetryError{}

	testResp := &http.Response{StatusCode: http.StatusNotFound}
	lre.Set(testResp, nil)

	err := lre.Get()
	// When err is nil but resp is non-nil, Set constructs a ResponseError
	if err == nil {
		t.Fatal("expected non-nil error (ResponseError constructed from response), got nil")
	}
	// The ResponseError should contain the status code
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to contain '404', got: %s", err.Error())
	}
}

func Test_LastRetryError_SetWithNilBoth(t *testing.T) {
	lre := &clients.LastRetryError{}
	lre.Set(nil, nil)

	err := lre.Get()
	if err != nil {
		t.Errorf("expected nil error when both resp and err are nil, got: %v", err)
	}
}

func Test_WrapContextError_NilError(t *testing.T) {
	lre := &clients.LastRetryError{}
	lre.Set(nil, errors.New("captured"))

	result := clients.WrapContextError(nil, lre)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func Test_WrapContextError_NilLastRetryError(t *testing.T) {
	err := context.DeadlineExceeded
	result := clients.WrapContextError(err, nil)
	if result != err {
		t.Errorf("expected %v, got %v", err, result)
	}
}

func Test_WrapContextError_NonContextError(t *testing.T) {
	lre := &clients.LastRetryError{}
	lre.Set(nil, errors.New("captured"))

	originalErr := errors.New("some other error")
	result := clients.WrapContextError(originalErr, lre)
	if result != originalErr {
		t.Errorf("expected original error %v, got %v", originalErr, result)
	}
}

func Test_WrapContextError_DeadlineExceeded_WithCapturedError(t *testing.T) {
	lre := &clients.LastRetryError{}
	capturedErr := errors.New("connection refused to api.example.com")
	lre.Set(nil, capturedErr)

	result := clients.WrapContextError(context.DeadlineExceeded, lre)

	// Should still be unwrappable to context.DeadlineExceeded
	if !errors.Is(result, context.DeadlineExceeded) {
		t.Errorf("expected error to wrap context.DeadlineExceeded")
	}

	// Should contain the captured error message
	expected := "context deadline exceeded, last retryable error: connection refused to api.example.com"
	if result.Error() != expected {
		t.Errorf("expected %q, got %q", expected, result.Error())
	}
}

func Test_WrapContextError_Canceled_WithCapturedError(t *testing.T) {
	lre := &clients.LastRetryError{}
	capturedErr := errors.New("server returned 500")
	lre.Set(nil, capturedErr)

	result := clients.WrapContextError(context.Canceled, lre)

	if !errors.Is(result, context.Canceled) {
		t.Errorf("expected error to wrap context.Canceled")
	}

	expected := "context canceled, last retryable error: server returned 500"
	if result.Error() != expected {
		t.Errorf("expected %q, got %q", expected, result.Error())
	}
}

func Test_WrapContextError_DeadlineExceeded_WithCapturedResponseError(t *testing.T) {
	lre := &clients.LastRetryError{}
	// Set with nil err and a response — should construct a ResponseError internally
	testResp := &http.Response{StatusCode: http.StatusTooManyRequests}
	lre.Set(testResp, nil)

	result := clients.WrapContextError(context.DeadlineExceeded, lre)

	if !errors.Is(result, context.DeadlineExceeded) {
		t.Errorf("expected error to wrap context.DeadlineExceeded")
	}

	// Should contain the status code from the ResponseError
	if !strings.Contains(result.Error(), "429") {
		t.Errorf("expected error to contain '429', got: %s", result.Error())
	}
	if !strings.Contains(result.Error(), "last retryable error:") {
		t.Errorf("expected error to contain 'last retryable error:', got: %s", result.Error())
	}
}

func Test_WrapContextError_DeadlineExceeded_WithNoCapturedInfo(t *testing.T) {
	lre := &clients.LastRetryError{}

	result := clients.WrapContextError(context.DeadlineExceeded, lre)

	// Should return original error unchanged
	if result != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", result)
	}
}

func Test_WrapContextError_DeadlineExceeded_ErrorTakesPrecedenceOverResponse(t *testing.T) {
	lre := &clients.LastRetryError{}
	capturedErr := errors.New("network timeout")
	testResp := &http.Response{StatusCode: http.StatusInternalServerError}
	lre.Set(testResp, capturedErr)

	result := clients.WrapContextError(context.DeadlineExceeded, lre)

	// When err is non-nil, it takes precedence over constructing from response
	expected := "context deadline exceeded, last retryable error: network timeout"
	if result.Error() != expected {
		t.Errorf("expected %q, got %q", expected, result.Error())
	}
}
