package clients

import (
	"context"
	"regexp"
	"slices"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// configureCustomRetry configures the retry configuration based on the supplied retry configuration.
// Using a dedicated function to allow for easier testing.
func configureCustomRetry(ctx context.Context, rtry retry.RetryValue, useReadAfterCreateValues bool) (*backoff.ExponentialBackOff, []regexp.Regexp, []int, []func(d interface{}) bool) {
	// Configure default retry configuration.
	// The default is to retry on 429 codes, so using the context deadline as max elapsed time is sane.
	maxElapsed := 2 * time.Minute
	if ctxDeadline, ok := ctx.Deadline(); ok {
		maxElapsed = time.Until(ctxDeadline)
	}
	backOff := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(5*time.Second),
		backoff.WithMaxInterval(30*time.Second),
		backoff.WithMaxElapsedTime(maxElapsed),
	)
	errRegExps := []regexp.Regexp{}
	statusCodes := rtry.GetDefaultRetryableStatusCodes()
	// Add default read after create values for the default retry configuration.
	if useReadAfterCreateValues {
		statusCodes = rtry.GetDefaultRetryableReadAfterCreateStatusCodes()
	}
	var dataCallbackFuncs []func(d interface{}) bool

	// If a custom retry configuration is supplied, use it.
	if !rtry.IsNull() && !rtry.IsUnknown() {
		tflog.Debug(ctx, "using custom retry configuration")

		statusCodes = rtry.GetRetryableStatusCodes()
		if useReadAfterCreateValues && rtry.ReadAfterCreateRetryOnNotFoundOrForbidden.ValueBool() {
			for _, code := range rtry.GetDefaultRetryableReadAfterCreateStatusCodes() {
				if !slices.Contains(statusCodes, code) {
					statusCodes = append(statusCodes, code)
				}
			}
		}
		if rtry.ResponseIsNil.ValueBool() {
			dataCallbackFuncs = []func(d interface{}) bool{
				func(d interface{}) bool {
					return d == nil
				},
			}
		}
		backOff = backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(rtry.GetIntervalSecondsAsDuration()),
			backoff.WithMaxInterval(rtry.GetMaxIntervalSecondsAsDuration()),
			backoff.WithMultiplier(rtry.GetMultiplier()),
			backoff.WithRandomizationFactor(rtry.GetRandomizationFactor()),
			backoff.WithMaxElapsedTime(maxElapsed),
		)
		errRegExps = rtry.GetErrorMessagesRegex()
	}

	return backOff, errRegExps, statusCodes, dataCallbackFuncs
}
