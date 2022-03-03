package clients

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	moduleName    = "resource"
	moduleVersion = "v0.1.0"
)

type ResourceClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

func NewResourceClient(subscriptionID string, credential azcore.TokenCredential, opt *arm.ClientOptions) *ResourceClient {
	if opt == nil {
		opt = &arm.ClientOptions{}
	}
	if opt.Endpoint == "" {
		opt.Endpoint = arm.AzurePublicCloud
	}
	return &ResourceClient{
		subscriptionID: subscriptionID,
		host:           string(opt.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, opt),
	}
}

func (client *ResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}) (interface{}, *http.Response, error) {
	resp, err := client.createOrUpdate(ctx, resourceID, apiVersion, body)
	if err != nil {
		return nil, nil, err
	}
	var responseBody interface{}
	pt, err := armruntime.NewPoller("Client.CreateOrUpdate", "", resp, client.pl)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, 10*time.Second, &responseBody)
		return responseBody, resp, err
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, nil, err
	}
	return responseBody, resp, nil
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
	urlPath := "/{resourceId}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
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

func (client *ResourceClient) Get(ctx context.Context, resourceID string, apiVersion string) (interface{}, *http.Response, error) {
	req, err := client.getCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, nil, runtime.NewResponseError(resp)
	}

	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, nil, err
	}
	return responseBody, resp, nil
}

func (client *ResourceClient) getCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := "/{resourceId}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
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

func (client *ResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string) (interface{}, *http.Response, error) {
	resp, err := client.delete(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, nil, err
	}
	var responseBody interface{}
	pt, err := armruntime.NewPoller("Client.Delete", "", resp, client.pl)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, 10*time.Second, &responseBody)
		return responseBody, resp, err
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, nil, err
	}
	return responseBody, resp, nil
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
	urlPath := "/{resourceId}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
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
