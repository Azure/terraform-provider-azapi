package clients_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/cenkalti/backoff/v4"
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
	t.Parallel()
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	rexs := clients.StringSliceToRegexpSliceMust([]string{"retry error"})
	bkof := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithMaxInterval(30*time.Second),
		backoff.WithMultiplier(2),
		backoff.WithRandomizationFactor(0.0),
	)
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, rexs, nil, nil)
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
	t.Parallel()
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	rexs := clients.StringSliceToRegexpSliceMust([]string{"^retry"})
	bkof := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMultiplier(1.5),
		backoff.WithRandomizationFactor(0.0),
	)
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, rexs, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientMultiRegexp(t *testing.T) {
	t.Parallel()
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	rexs := clients.StringSliceToRegexpSliceMust([]string{"nomatch", "^retry"})
	bkof := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMultiplier(1.5),
		backoff.WithRandomizationFactor(0.0),
	)
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, rexs, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientMultiRegexpNoMatchWithPermError(t *testing.T) {
	t.Parallel()
	mock := NewMockResourceClient(t, nil, errors.New("perm error"), 3, errors.New("retry error"))
	rexs := clients.StringSliceToRegexpSliceMust([]string{"retry"})
	bkof := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMultiplier(1.5),
		backoff.WithRandomizationFactor(0.0),
	)
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, rexs, nil, nil)
	_, err := retryClient.Get(context.Background(), "", "", clients.DefaultRequestOptions())
	assert.ErrorContains(t, err, "perm error")
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientContextDeadline(t *testing.T) {
	t.Parallel()
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(60*time.Second),
		backoff.WithMaxInterval(60*time.Second),
		backoff.WithMultiplier(1.5),
		backoff.WithRandomizationFactor(0.0),
	)
	rexs := clients.StringSliceToRegexpSliceMust([]string{"^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, rexs, nil, nil)
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
