package utils

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func ResponseErrorWasNotFound(err error) bool {
	return ResponseErrorWasStatusCode(err, http.StatusNotFound)
}

func ResponseErrorWasStatusCode(err error, statusCode int) bool {
	var responseErr *azcore.ResponseError
	return errors.As(err, &responseErr) && responseErr.StatusCode == statusCode
}
