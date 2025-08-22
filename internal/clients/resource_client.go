package clients

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	moduleName    = "resource"
	moduleVersion = "v0.1.0"
)

type ResourceClient struct {
	host string
	pl   runtime.Pipeline
}

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

func (client *ResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}, options RequestOptions) (interface{}, error) {
	// override the default retry options with the ones provided in the options
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)

		log.Printf("[DEBUG] Retry configuration is custom: MaxRetries %d, RetryDelay %v, MaxRetryDelay %v, StatusCodes %v, ShouldRetryFunc %t",
			options.RetryOptions.MaxRetries,
			options.RetryOptions.RetryDelay,
			options.RetryOptions.MaxRetryDelay,
			options.RetryOptions.StatusCodes,
			options.RetryOptions.ShouldRetry != nil,
		)
	}

	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}

	// Set the query parameters
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	// Set the headers
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}
	err = runtime.MarshalAsJSON(req, body)
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

func (client *ResourceClient) Get(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error) {
	// override the default retry options with the ones provided in the options
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)
		log.Printf("[DEBUG] Retry configuration is custom: MaxRetries %d, RetryDelay %v, MaxRetryDelay %v, StatusCodes %v, ShouldRetryFunc %t",
			options.RetryOptions.MaxRetries,
			options.RetryOptions.RetryDelay,
			options.RetryOptions.MaxRetryDelay,
			options.RetryOptions.StatusCodes,
			options.RetryOptions.ShouldRetry != nil,
		)
	}

	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}

	// Set the query parameters
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	// Set the headers
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
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

func (client *ResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string, options RequestOptions) (interface{}, error) {
	// override the default retry options with the ones provided in the options
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)
		log.Printf("[DEBUG] Retry configuration is custom: MaxRetries %d, RetryDelay %v, MaxRetryDelay %v, StatusCodes %v, ShouldRetryFunc %t",
			options.RetryOptions.MaxRetries,
			options.RetryOptions.RetryDelay,
			options.RetryOptions.MaxRetryDelay,
			options.RetryOptions.StatusCodes,
			options.RetryOptions.ShouldRetry != nil,
		)
	}

	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}

	// Set the query parameters
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	// Set the headers
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}

	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
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

func (client *ResourceClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (interface{}, error) {
	// override the default retry options with the ones provided in the options
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)
		log.Printf("[DEBUG] Retry configuration is custom: MaxRetries %d, RetryDelay %v, MaxRetryDelay %v, StatusCodes %v, ShouldRetryFunc %t",
			options.RetryOptions.MaxRetries,
			options.RetryOptions.RetryDelay,
			options.RetryOptions.MaxRetryDelay,
			options.RetryOptions.StatusCodes,
			options.RetryOptions.ShouldRetry != nil,
		)
	}
	urlPath := resourceID
	if action != "" {
		urlPath = fmt.Sprintf("%s/%s", resourceID, action)
	}
	req, err := runtime.NewRequest(ctx, method, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}

	// Set the query parameters
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	// Set the headers
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}

	// Set the body if method is not GET
	if method != "GET" && body != nil {
		if err = runtime.MarshalAsJSON(req, body); err != nil {
			return nil, err
		}
	}

	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
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

func (client *ResourceClient) List(ctx context.Context, url string, apiVersion string, options RequestOptions) (interface{}, error) {
	// override the default retry options with the ones provided in the options
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)
		log.Printf("[DEBUG] Retry configuration is custom: MaxRetries %d, RetryDelay %v, MaxRetryDelay %v, StatusCodes %v, ShouldRetryFunc %t",
			options.RetryOptions.MaxRetries,
			options.RetryOptions.RetryDelay,
			options.RetryOptions.MaxRetryDelay,
			options.RetryOptions.StatusCodes,
			options.RetryOptions.ShouldRetry != nil,
		)
	}

	pager := runtime.NewPager(runtime.PagingHandler[interface{}]{
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
