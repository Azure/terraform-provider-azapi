package clients_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/stretchr/testify/assert"
)

// TestRetryClient tests the retry client with a mock client that returns an error.
// This test should retry 3 times (4 requests in total) and retirn a nil error.
// The ransomized interval is calculated by the backoff package as follows:
//
//	randomized interval =
//	  RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])
//
// Therefore a configuration of interval = 1s, multiplier = 2 and randomization = 0.0 should result in
// the following intervals:
// 1. 0s (the 1st request is made immediately)
// 2. 1s (iniital interval)
// 3. 2s (last interval 1s * 2 multiplier)
// 4. 4s (last interval 2s * 2 multiplier)
//
// We check these timings are expected using the assert.InDeltaSlice function.
func TestRetryClient(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(1, 30, 2, 0.0, []string{"retry error"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.requestCount)
	assert.Len(t, mock.requestTimes, 4)
	requestBackoffDiractionNanos := make([]int, len(mock.requestTimes))
	for i, v := range mock.requestTimes {
		if i == 0 {
			requestBackoffDiractionNanos[i] = 0
			continue
		}
		requestBackoffDiractionNanos[i] = int(v.Sub(mock.requestTimes[i-1]).Nanoseconds())
	}
	// We assert that the intervals are within 0.01 seconds of the expected values
	assert.InDeltaSlice(t, []int{0, 1e+09, 2e+09, 4e+09}, requestBackoffDiractionNanos, 1e+07)
}

func TestRetryClientRegexp(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(1, 5, 1.5, 0.0, []string{"^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientMultiRegexp(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(1, 5, 1.5, 0.0, []string{"nomatch", "^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientMultiRegexpNoMatchWithPermError(t *testing.T) {
	mock := NewMockResourceClient(t, nil, errors.New("perm error"), 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(1, 5, 1.5, 0.0, []string{"retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.ErrorContains(t, err, "perm error")
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientContextDeadline(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(60, 60, 1.5, 0.0, []string{"^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps, nil, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()
	_, err := retryClient.Get(ctx, "", "", clients.DefaultRequestOptions())
	elapsed := time.Since(start)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.True(t, elapsed < 15*time.Second)
	// Test that the context was cancelled
	_, ok := <-ctx.Done()
	assert.False(t, ok)
}
