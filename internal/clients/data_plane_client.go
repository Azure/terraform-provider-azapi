package clients

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type RawRequestOptions struct {
	RequestOptions
	APIVersionHeader string
	Content          []byte
	SuccessCodes     []int
}

type DataPlaneClient struct {
	credential      azcore.TokenCredential
	clientOptions   *arm.ClientOptions
	cachedPipelines map[string]runtime.Pipeline
	syncMux         sync.Mutex
}

func NewDataPlaneClient(credential azcore.TokenCredential, opt *arm.ClientOptions) (*DataPlaneClient, error) {
	if opt == nil {
		opt = &arm.ClientOptions{}
	}
	return &DataPlaneClient{
		credential:      credential,
		clientOptions:   opt,
		cachedPipelines: make(map[string]runtime.Pipeline),
		syncMux:         sync.Mutex{},
	}, nil
}

func (client *DataPlaneClient) CreateOrUpdateThenPoll(ctx context.Context, id parse.DataPlaneResourceId, body interface{}, options RequestOptions) (interface{}, error) {
	urlPath := buildURL(id.AzureResourceId, "")

	req, err := buildRequest(ctx, options, urlPath, http.MethodPut, id.ApiVersion)
	if err != nil {
		return nil, err
	}

	err = runtime.MarshalAsJSON(req, body)
	if err != nil {
		return nil, err
	}

	successCodes := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent}
	return client.sendRequestThenPoll(ctx, req, urlPath, options, successCodes)
}

func (client *DataPlaneClient) Get(ctx context.Context, id parse.DataPlaneResourceId, options RequestOptions) (interface{}, error) {
	urlPath := buildURL(id.AzureResourceId, "")
	req, err := buildRequest(ctx, options, urlPath, http.MethodGet, id.ApiVersion)
	if err != nil {
		return nil, err
	}

	successCodes := []int{http.StatusOK}
	resp, _, err := client.sendRequest(req, urlPath, options, successCodes)
	if err != nil {
		return nil, err
	}

	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *DataPlaneClient) DeleteThenPoll(ctx context.Context, id parse.DataPlaneResourceId, options RequestOptions) (interface{}, error) {
	urlPath := buildURL(id.AzureResourceId, "")
	req, err := buildRequest(ctx, options, urlPath, http.MethodDelete, id.ApiVersion)
	if err != nil {
		return nil, err
	}

	successCodes := []int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}
	return client.sendRequestThenPoll(ctx, req, urlPath, options, successCodes)
}

func (client *DataPlaneClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options RequestOptions) (interface{}, error) {
	urlPath := buildURL(resourceID, action)
	req, err := buildRequest(ctx, options, urlPath, method, apiVersion)
	if err != nil {
		return nil, err
	}

	if method != "GET" && body != nil {
		err = runtime.MarshalAsJSON(req, body)
	}
	if err != nil {
		return nil, err
	}

	// Action does not use sendRequestThenPoll because it parses the response body
	// based on Content-Type (text/plain vs application/json) rather than always as JSON.
	successCodes := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}
	resp, pipeline, err := client.sendRequest(req, urlPath, options, successCodes)
	if err != nil {
		return nil, err
	}

	pt, err := runtime.NewPoller[interface{}](resp, pipeline, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}

	var responseBody interface{}
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

func (client *DataPlaneClient) DoRaw(ctx context.Context, method string, resourceID string, options RawRequestOptions) (*http.Response, error) {
	urlPath, err := buildEscapedURL(resourceID)
	if err != nil {
		return nil, err
	}

	req, err := buildRawRequest(ctx, options, urlPath, method)
	if err != nil {
		return nil, err
	}
	if options.Content != nil {
		if err := req.SetBody(streaming.NopCloser(bytes.NewReader(options.Content)), ""); err != nil {
			return nil, err
		}
		req.Raw().ContentLength = int64(len(options.Content))
		req.Raw().Header.Set("Content-Length", fmt.Sprintf("%d", len(options.Content)))
	}
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}
	if options.APIVersionHeader != "" {
		req.Raw().Header.Set("x-ms-version", options.APIVersionHeader)
	}

	successCodes := options.SuccessCodes
	if len(successCodes) == 0 {
		successCodes = []int{http.StatusOK}
	}
	resp, _, err := client.sendRequestWithServiceRequirement(req, urlPath, options.RequestOptions, successCodes, options.APIVersionHeader != "")
	return resp, err
}

func buildURL(resourceId string, action string) (urlPath string) {
	urlPath = fmt.Sprintf("https://%s", resourceId)
	if action != "" {
		urlPath = fmt.Sprintf("%s/%s", urlPath, action)
	}
	return
}

func buildEscapedURL(resourceID string) (string, error) {
	host, path, found := strings.Cut(resourceID, "/")
	if host == "" {
		return "", fmt.Errorf("invalid data-plane resource ID %q: missing host", resourceID)
	}
	if !found {
		return (&url.URL{Scheme: "https", Host: host}).String(), nil
	}
	return (&url.URL{Scheme: "https", Host: host, Path: "/" + path}).String(), nil
}

func buildRawRequest(ctx context.Context, options RawRequestOptions, urlPath, method string) (*policy.Request, error) {
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)
	}
	req, err := runtime.NewRequest(ctx, method, urlPath)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

func buildRequest(ctx context.Context, options RequestOptions, urlPath, method, apiVersion string) (*policy.Request, error) {
	if options.RetryOptions != nil {
		ctx = policy.WithRetryOptions(ctx, *options.RetryOptions)
	}
	req, err := runtime.NewRequest(ctx, method, urlPath)
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

func (client *DataPlaneClient) sendRequest(req *policy.Request, urlPath string, options RequestOptions, statusCodes []int) (*http.Response, runtime.Pipeline, error) {
	return client.sendRequestWithServiceRequirement(req, urlPath, options, statusCodes, false)
}

func (client *DataPlaneClient) sendRequestWithServiceRequirement(req *policy.Request, urlPath string, options RequestOptions, statusCodes []int, requireDataPlaneService bool) (*http.Response, runtime.Pipeline, error) {
	pipeline, err := client.cachedPipeline(urlPath, requireDataPlaneService)
	if err != nil {
		return nil, runtime.Pipeline{}, err
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		return nil, runtime.Pipeline{}, WrapContextError(err, options.LastRetryError)
	}
	if !runtime.HasStatusCode(resp, statusCodes...) {
		return nil, runtime.Pipeline{}, runtime.NewResponseError(resp)
	}
	return resp, pipeline, nil
}

func (client *DataPlaneClient) sendRequestThenPoll(ctx context.Context, req *policy.Request, urlPath string, options RequestOptions, statusCodes []int) (interface{}, error) {
	resp, pipeline, err := client.sendRequest(req, urlPath, options, statusCodes)
	if err != nil {
		return nil, err
	}

	if pt, err := runtime.NewPoller[interface{}](resp, pipeline, nil); err == nil {
		return pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
	}

	// if NewPoller returned an error, return the original resp body directly.
	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *DataPlaneClient) cachedPipeline(rawUrl string, requireDataPlaneService bool) (runtime.Pipeline, error) {
	client.syncMux.Lock()
	defer client.syncMux.Unlock()

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return runtime.Pipeline{}, err
	}
	serviceName := cloud.ResourceManager
	matchedService := false
	cloudConfig := client.clientOptions.Cloud
	host := parsedUrl.Host
	for name, serviceConfiguration := range cloudConfig.Services {
		if strings.HasSuffix(host, strings.TrimPrefix(serviceConfiguration.Endpoint, "https://")) {
			serviceName = name
			matchedService = true
			break
		}
	}
	if requireDataPlaneService && (!matchedService || serviceName == cloud.ResourceManager) {
		return runtime.Pipeline{}, fmt.Errorf("unsupported data-plane endpoint %q: no matching service endpoint and token audience are configured for this cloud", host)
	}

	if pipeline, ok := client.cachedPipelines[string(serviceName)]; ok {
		return pipeline, nil
	}

	plOpt := runtime.PipelineOptions{}
	plOpt.APIVersion.Name = "api-version"
	authPolicy := armruntime.NewBearerTokenPolicy(client.credential, &armpolicy.BearerTokenOptions{Scopes: []string{cloudConfig.Services[serviceName].Audience + "/.default"}})
	plOpt.PerRetry = append(plOpt.PerRetry, authPolicy)
	pl := runtime.NewPipeline(moduleName, moduleVersion, plOpt, &client.clientOptions.ClientOptions)

	client.cachedPipelines[string(serviceName)] = pl
	return pl, nil
}
