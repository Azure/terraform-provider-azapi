package utils

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func ResponseWasForbidden(err error) bool {
	return ResponseErrorWasStatusCode(err, http.StatusForbidden)
}

func ResponseErrorWasNotFound(err error) bool {
	return ResponseErrorWasStatusCode(err, http.StatusNotFound)
}

// ResponseErrorWasNoRegisteredProvider checks if the error is an HTTP 400
// with error code "NoRegisteredProviderFound". This typically happens when the
// embedded API version in the provider is not (or no longer) supported by Azure,
// e.g. a preview version that has been replaced by the GA release.
func ResponseErrorWasNoRegisteredProvider(err error) bool {
	var responseErr *azcore.ResponseError
	return errors.As(err, &responseErr) &&
		responseErr.StatusCode == http.StatusBadRequest &&
		responseErr.ErrorCode == "NoRegisteredProviderFound"
}

func ResponseErrorWasStatusCode(err error, statusCode int) bool {
	var responseErr *azcore.ResponseError
	return errors.As(err, &responseErr) && responseErr.StatusCode == statusCode
}
