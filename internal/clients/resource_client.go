package clients

import (
	"context"
	"fmt"
	"net/http"
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
	host  string
	pl    runtime.Pipeline
	retry *ResourceClientRetryableErrors
}

type ResourceClientRetryableErrors struct {
	backoff *backoff.ExponentialBackOff
	errors  []string
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
		host:  ep,
		pl:    pl,
		retry: nil,
	}, nil
}

// NewResourceClientRertryableErrors creates a ResourceClientRetryableErrors object.
// TODO: add inputs with resource schema to generate the retryable errors object.
func NewResourceClientRetryableErrors() *ResourceClientRetryableErrors {
	return &ResourceClientRetryableErrors{}
}

// WithRetry configures the retryable errors for the client.
func (client *ResourceClient) WithRetry(retry *ResourceClientRetryableErrors) *ResourceClient {
	client.retry = retry
	return client
}

// CreateOrUpdateWithRetry configures the retryable errors for the client.
// It calls CreateOrUpdate, then checks if the error is contained in the retryable errors list.
// If it is, it will retry the operation with the configured backoff.
// If it is not, it will return the error as a backoff.PermanentError{}.
func (client *ResourceClient) CreateOrUpdateWithRetry(ctx context.Context, resourceID string, apiVersion string, body interface{}) (interface{}, error) {
	if client.retry == nil {
		return nil, fmt.Errorf("Retry is not configured. Please call WithRetry() first.")
	}
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := client.CreateOrUpdate(ctx, resourceID, apiVersion, body)
			if err != nil {
				for _, e := range client.retry.errors {
					if !strings.Contains(err.Error(), e) {
						return nil, &backoff.PermanentError{Err: err}
					}
				}
			}
			return data, err
		})
	exbo := backoff.WithContext(client.retry.backoff, ctx)
	return backoff.RetryWithData[interface{}](op, exbo)
}

func (client *ResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}) (interface{}, error) {
	resp, err := client.createOrUpdate(ctx, resourceID, apiVersion, body)
	if err != nil {
		return nil, err
	}
	var responseBody interface{}
	pt, err := runtime.NewPoller[interface{}](resp, client.pl, nil)
	if err != nil {
		return nil, err
	}
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
		return ptresp, nil
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *ResourceClient) createOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceID, apiVersion, body)
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

func (client *ResourceClient) createOrUpdateCreateRequest(ctx context.Context, resourceID string, apiVersion string, body interface{}) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, body)
}

func (client *ResourceClient) Get(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	req, err := client.getCreateRequest(ctx, resourceID, apiVersion)
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

func (client *ResourceClient) getCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

func (client *ResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	resp, err := client.delete(ctx, resourceID, apiVersion)
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

func (client *ResourceClient) delete(ctx context.Context, resourceID string, apiVersion string) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceID, apiVersion)
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

func (client *ResourceClient) deleteCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

func (client *ResourceClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (interface{}, error) {
	resp, err := client.action(ctx, resourceID, action, apiVersion, method, body)
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

func (client *ResourceClient) action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (*http.Response, error) {
	req, err := client.actionCreateRequest(ctx, resourceID, action, apiVersion, method, body)
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

func (client *ResourceClient) actionCreateRequest(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (*policy.Request, error) {
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
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	if method != "GET" && body != nil {
		err = runtime.MarshalAsJSON(req, body)
	}
	return req, err
}

func (client *ResourceClient) List(ctx context.Context, url string, apiVersion string) (interface{}, error) {
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
				req.Raw().URL.RawQuery = reqQP.Encode()
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
	if responseErr, ok := err.(*azcore.ResponseError); ok {
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
