package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/cenkalti/backoff/v4"
)

const (
	moduleName    = "resource"
	moduleVersion = "v0.1.0"
)

type ResourceClient struct {
	host string
	pl   runtime.Pipeline
}

// ResourceClientRetryableErrors is a wrapper around ResourceClient that allows for retrying on specific errors.
type ResourceClientRetryableErrors struct {
	client            Requester                   // client is a Requester interface to allow mocking
	backoff           *backoff.ExponentialBackOff // backoff is the backoff configuration for retrying
	errors            []regexp.Regexp             // errors is the list of errors regexp to retry on
	statusCodes       []int                       // statusCodes is the list of status codes to retry on
	dataCallbackFuncs []func(interface{}) bool    // dataCallbackFuncs is the list of functions to call to determine if the data is retryable
}

// NewResourceClientRetryableErrors creates a new ResourceClientRetryableErrors.
func NewResourceClientRetryableErrors(client Requester, bkof *backoff.ExponentialBackOff, errRegExps []regexp.Regexp, statusCodes []int, dataCallbackFuncs []func(any) bool) *ResourceClientRetryableErrors {
	rcre := &ResourceClientRetryableErrors{
		client:            client,
		backoff:           bkof,
		errors:            errRegExps,
		statusCodes:       statusCodes,
		dataCallbackFuncs: dataCallbackFuncs,
	}
	rcre.backoff.Reset()
	return rcre
}

// Requester is the interface for HTTP operations, meaning we can supply a ResourceClient or a ResourceClientRetryableErrors.
type Requester interface {
	Get(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error)
	CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}, options RequestOptions) (interface{}, error)
	Delete(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error)
	Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (interface{}, error)
	List(ctx context.Context, url string, apiVersion string, options RequestOptions) (interface{}, error)
}

var (
	_ Requester = &ResourceClient{}
	_ Requester = &ResourceClientRetryableErrors{}
)

func NewResourceClient(credential azcore.TokenCredential, opt *arm.ClientOptions) (*ResourceClient, error) {
	if opt == nil {
		opt = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := opt.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, opt)
	if err != nil {
		return nil, err
	}
	return &ResourceClient{
		host: ep,
		pl:   pl,
	}, nil
}

// NewRetryableErrors creates the backoff and error regexs for retryable errors.
func NewRetryableErrors(intervalSeconds, maxIntervalSeconds int, multiplier, randomizationFactor float64, errorRegexs []string) (*backoff.ExponentialBackOff, []regexp.Regexp) {
	bkof := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(time.Duration(intervalSeconds)*time.Second),
		backoff.WithRandomizationFactor(randomizationFactor),
		backoff.WithMaxInterval(time.Duration(maxIntervalSeconds)*time.Second),
		backoff.WithRandomizationFactor(randomizationFactor),
		backoff.WithMultiplier(multiplier),
	)
	res := make([]regexp.Regexp, len(errorRegexs))
	for i, e := range errorRegexs {
		res[i] = *regexp.MustCompile(e) // MustCompile as schema has custom validation so we know it's valid
	}
	return bkof, res
}

// WithRetry configures the retryable errors for the client.
func (client *ResourceClient) WithRetry(bkof *backoff.ExponentialBackOff, errRegExps []regexp.Regexp, statusCodes []int, dataCallbackFuncs []func(interface{}) bool) *ResourceClientRetryableErrors {
	rcre := &ResourceClientRetryableErrors{
		client:            client,
		backoff:           bkof,
		errors:            errRegExps,
		statusCodes:       statusCodes,
		dataCallbackFuncs: dataCallbackFuncs,
	}
	rcre.backoff.Reset()
	return rcre
}

// CreateOrUpdate configures the retryable errors for the client.
// It calls CreateOrUpdate, then checks if the error is contained in the retryable errors list.
// If it is, it will retry the operation with the configured backoff.
// If it is not, it will return the error as a backoff.PermanentError{}.
func (retryclient *ResourceClientRetryableErrors) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.CreateOrUpdate(ctx, resourceID, apiVersion, body, options)
			if err != nil {
				if isRetryable(*retryclient, data, err) {
					return data, err
				}
				return nil, &backoff.PermanentError{Err: err}
			}
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData[interface{}](op, exbo)
}

func (client *ResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}, options RequestOptions) (interface{}, error) {
	resp, err := client.createOrUpdate(ctx, resourceID, apiVersion, body, options)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		if err == nil {
			return resp, nil
		}
		if !client.shouldIgnorePollingError(err) {
			return nil, err
		}
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) createOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}, options RequestOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceID, apiVersion, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func (client *ResourceClient) createOrUpdateCreateRequest(ctx context.Context, resourceID string, apiVersion string, body interface{}, options RequestOptions) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}
	return req, runtime.MarshalAsJSON(req, body)
}

// Get configures the retryable errors for the client.
// It calls Get, then checks if the error is contained in the retryable errors list.
// If it is, it will retry the operation with the configured backoff.
// If it is not, it will return the error as a backoff.PermanentError{}.
func (retryclient *ResourceClientRetryableErrors) Get(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.Get(ctx, resourceID, apiVersion, options)
			if err != nil {
				if isRetryable(*retryclient, data, err) {
					return data, err
				}
				return nil, &backoff.PermanentError{Err: err}
			}
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData[interface{}](op, exbo)
}

func (client *ResourceClient) Get(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error) {
	req, err := client.getCreateRequest(ctx, resourceID, apiVersion, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}

	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) getCreateRequest(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}
	return req, nil
}

// Delete configures the retryable errors for the client.
// It calls Delete, then checks if the error is contained in the retryable errors list.
// If it is, it will retry the operation with the configured backoff.
// If it is not, it will return the error as a backoff.PermanentError{}.
func (retryclient *ResourceClientRetryableErrors) Delete(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.Delete(ctx, resourceID, apiVersion, options)
			if err != nil {
				if isRetryable(*retryclient, data, err) {
					return data, err
				}
				return nil, &backoff.PermanentError{Err: err}
			}
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData[interface{}](op, exbo)
}

func (client *ResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error) {
	resp, err := client.delete(ctx, resourceID, apiVersion, options)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		if err == nil {
			return resp, nil
		}
		if !client.shouldIgnorePollingError(err) {
			return nil, err
		}
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) delete(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceID, apiVersion, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func (client *ResourceClient) deleteCreateRequest(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}
	return req, nil
}

// Action configures the retryable errors for the client.
// It calls Action, then checks if the error is contained in the retryable errors list.
// If it is, it will retry the operation with the configured backoff.
// If it is not, it will return the error as a backoff.PermanentError{}.
func (retryclient *ResourceClientRetryableErrors) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.Action(ctx, resourceID, action, apiVersion, method, body, options)
			if err != nil {
				if isRetryable(*retryclient, data, err) {
					return data, err
				}
				return nil, &backoff.PermanentError{Err: err}
			}
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData[interface{}](op, exbo)
}

func (client *ResourceClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (interface{}, error) {
	resp, err := client.action(ctx, resourceID, action, apiVersion, method, body, options)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		if err == nil {
			return resp, nil
		}
		if !client.shouldIgnorePollingError(err) {
			return nil, err
		}
	}

	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "text/plain"):
		payload, err := runtime.Payload(resp)
		if err != nil {
			return nil, err
		}
		responseBody = string(payload)
	case strings.Contains(contentType, "application/json"):
		if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
			return nil, err
		}
	default:
	}
	return responseBody, nil
}

func (client *ResourceClient) action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (*http.Response, error) {
	req, err := client.actionCreateRequest(ctx, resourceID, action, apiVersion, method, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func (client *ResourceClient) actionCreateRequest(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (*policy.Request, error) {
	urlPath := resourceID
	if action != "" {
		urlPath = fmt.Sprintf("%s/%s", resourceID, action)
	}
	req, err := runtime.NewRequest(ctx, method, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}
	if method != "GET" && body != nil {
		err = runtime.MarshalAsJSON(req, body)
	}
	return req, err
}

// List configures the retryable errors for the client.
// It calls Get, then checks if the error is contained in the retryable errors list.
// If it is, it will retry the operation with the configured backoff.
// If it is not, it will return the error as a backoff.PermanentError{}.
func (retryclient *ResourceClientRetryableErrors) List(ctx context.Context, url string, apiVersion string, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.List(ctx, url, apiVersion, options)
			if err != nil {
				if isRetryable(*retryclient, data, err) {
					return data, err
				}
				return nil, &backoff.PermanentError{Err: err}
			}
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData[interface{}](op, exbo)
}

func (client *ResourceClient) List(ctx context.Context, url string, apiVersion string, options RequestOptions) (interface{}, error) {
	pager := runtime.NewPager[interface{}](runtime.PagingHandler[interface{}]{
		More: func(current interface{}) bool {
			if current == nil {
				return false
			}
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return false
			}
			if currentMap["nextLink"] == nil {
				return false
			}
			if nextLink := currentMap["nextLink"].(string); nextLink == "" {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, current *interface{}) (interface{}, error) {
			var request *policy.Request
			if current == nil {
				req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, url))
				if err != nil {
					return nil, err
				}
				reqQP := req.Raw().URL.Query()
				reqQP.Set("api-version", apiVersion)
				for key, value := range options.QueryParameters {
					reqQP.Set(key, value)
				}
				req.Raw().URL.RawQuery = reqQP.Encode()
				for key, value := range options.Headers {
					req.Raw().Header.Set(key, value)
				}
				request = req
			} else {
				nextLink := ""
				if currentMap, ok := (*current).(map[string]interface{}); ok && currentMap["nextLink"] != nil {
					nextLink = currentMap["nextLink"].(string)
				}
				req, err := runtime.NewRequest(ctx, http.MethodGet, nextLink)
				if err != nil {
					return nil, err
				}
				request = req
			}
			request.Raw().Header.Set("Accept", "application/json")
			resp, err := client.pl.Do(request)
			if err != nil {
				return nil, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return nil, runtime.NewResponseError(resp)
			}
			var responseBody interface{}
			if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
				return nil, err
			}
			return responseBody, nil
		},
	})

	value := make([]interface{}, 0)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if pageMap, ok := page.(map[string]interface{}); ok {
			if pageMap["value"] != nil {
				if pageValue, ok := pageMap["value"].([]interface{}); ok {
					value = append(value, pageValue...)
					continue
				}
			}
		}

		// if response doesn't follow the ARM paging guideline, return the response as is
		return page, nil
	}
	return map[string]interface{}{
		"value": value,
	}, nil
}

func (client *ResourceClient) shouldIgnorePollingError(err error) bool {
	if err == nil {
		return true
	}
	// there are some APIs that don't follow the ARM LRO guideline, return the response as is
	var responseErr *azcore.ResponseError
	if errors.As(err, &responseErr) {
		if responseErr.RawResponse != nil && responseErr.RawResponse.Request != nil {
			// all control plane APIs must flow through ARM, ignore the polling error if it's not ARM
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/25356, in this case, the polling url is not exposed by ARM
			pollRequest := responseErr.RawResponse.Request
			if pollRequest.Host != strings.TrimPrefix(client.host, "https://") {
				return true
			}

			// ignore the polling error if the polling url doesn't support GET method
			// issue:https://github.com/Azure/azure-rest-api-specs/issues/25362, in this case, the polling url doesn't support GET method
			if responseErr.StatusCode == http.StatusMethodNotAllowed {
				return true
			}
		}
	}
	return false
}

func isRetryable(retryclient ResourceClientRetryableErrors, data interface{}, err error) bool {
	for _, e := range retryclient.errors {
		if e.MatchString(err.Error()) {
			return true
		}
	}
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		if slices.Contains(retryclient.statusCodes, respErr.StatusCode) {
			return true
		}
	}
	for _, f := range retryclient.dataCallbackFuncs {
		if f(data) {
			return true
		}
	}
	return false
}
