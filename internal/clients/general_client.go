package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
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
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultScopeSuffix  = "/.default"
	apiVersionHeaderKey = "api-version"
)

var (
	_ GeneralRequester = &GeneralClient{}
	_ GeneralRequester = &GeneralClientRetryableErrors{}
)

type GeneralRequester interface {
	CreateOrUpdateThenPoll(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, body any, options RequestOptions) (interface{}, error)
	Get(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, options RequestOptions) (interface{}, error)
	DeleteThenPoll(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, options RequestOptions) (interface{}, error)
	Action(ctx context.Context, url, apiVersion, method string, svcConfig cloud.ServiceConfiguration, body any, options RequestOptions) (interface{}, error)
}

type GeneralClientRetryableErrors struct {
	client            GeneralRequester            // client is a Requester interface to allow mocking
	backoff           *backoff.ExponentialBackOff // backoff is the backoff configuration for retrying
	errors            []regexp.Regexp             // errors is the list of errors regexp to retry on
	statusCodes       []int                       // statusCodes is the list of status codes to retry on
	dataCallbackFuncs []func(interface{}) bool    // dataCallbackFuncs is the list of functions to call to determine if the data is retryable
}

type GeneralClient struct {
	credential      azcore.TokenCredential
	clientOptions   *arm.ClientOptions
	cachedPipelines map[cloud.ServiceConfiguration]runtime.Pipeline
	syncMux         sync.Mutex
}

func NewGeneralClient(credential azcore.TokenCredential, options *arm.ClientOptions) *GeneralClient {
	if options == nil {
		options = &arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloud.AzurePublic,
			},
		}
	}

	return &GeneralClient{
		credential:      credential,
		clientOptions:   options,
		cachedPipelines: make(map[cloud.ServiceConfiguration]runtime.Pipeline),
		syncMux:         sync.Mutex{},
	}
}

// WithRetry configures the retryable errors for the client.
func (client *GeneralClient) WithRetry(bkof *backoff.ExponentialBackOff, errRegExps []regexp.Regexp, statusCodes []int, dataCallbackFuncs []func(interface{}) bool) *GeneralClientRetryableErrors {
	rcre := &GeneralClientRetryableErrors{
		client:            client,
		backoff:           bkof,
		errors:            errRegExps,
		statusCodes:       statusCodes,
		dataCallbackFuncs: dataCallbackFuncs,
	}
	rcre.backoff.Reset()
	return rcre
}

// NewGeneralClientRetryableErrors creates a new ResourceClientRetryableErrors.
func NewGeneralClientRetryableErrors(client GeneralRequester, bkof *backoff.ExponentialBackOff, errRegExps []regexp.Regexp, statusCodes []int, dataCallbackFuncs []func(any) bool) *GeneralClientRetryableErrors {
	rcre := &GeneralClientRetryableErrors{
		client:            client,
		backoff:           bkof,
		errors:            errRegExps,
		statusCodes:       statusCodes,
		dataCallbackFuncs: dataCallbackFuncs,
	}
	rcre.backoff.Reset()
	return rcre
}

func (client *GeneralClient) cachedPipeline(serviceConfig cloud.ServiceConfiguration) (runtime.Pipeline, error) {
	client.syncMux.Lock()
	defer client.syncMux.Unlock()

	if pipeline, ok := client.cachedPipelines[serviceConfig]; ok {
		return pipeline, nil
	}

	plOpt := runtime.PipelineOptions{}
	plOpt.APIVersion.Name = apiVersionHeaderKey
	bearerTokenOptions := &armpolicy.BearerTokenOptions{
		Scopes: []string{
			serviceConfig.Audience + defaultScopeSuffix,
		},
	}
	authPolicy := armruntime.NewBearerTokenPolicy(client.credential, bearerTokenOptions)
	plOpt.PerRetry = append(plOpt.PerRetry, authPolicy)
	pl := runtime.NewPipeline(moduleName, moduleVersion, plOpt, &client.clientOptions.ClientOptions)

	client.cachedPipelines[serviceConfig] = pl
	return pl, nil
}

func (retryclient *GeneralClientRetryableErrors) updateContext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, "backoff_max_elapsed_time", retryclient.backoff.MaxElapsedTime.String())
	ctx = tflog.SetField(ctx, "backoff_initial_interval", retryclient.backoff.InitialInterval.String())
	ctx = tflog.SetField(ctx, "backoff_max_interval", retryclient.backoff.MaxInterval.String())
	ctx = tflog.SetField(ctx, "backoff_multiplier", retryclient.backoff.Multiplier)
	ctx = tflog.SetField(ctx, "backoff_randomization_factor", retryclient.backoff.RandomizationFactor)
	ctx = tflog.SetField(ctx, "retryable_http_status_codes", retryclient.statusCodes)
	ctx = tflog.SetField(ctx, "retryable_data_callback_funcs_length", len(retryclient.dataCallbackFuncs))
	re := make([]string, len(retryclient.errors))
	for i, r := range retryclient.errors {
		re[i] = r.String()
	}
	ctx = tflog.SetField(ctx, "retryable_errors", re)
	return ctx
}

func (client *GeneralClient) Get(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, options RequestOptions) (interface{}, error) {
	// build request
	urlPath := fmt.Sprintf("https://%s", url)
	req, err := runtime.NewRequest(ctx, http.MethodGet, urlPath)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if apiVersion != "" {
		reqQP.Set(apiVersionHeaderKey, apiVersion)
	}

	// Add query parameters and headers from the inputs.
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}

	// Set the request URL and headers, accept JSON by default.
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")

	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}

	// send request
	pipeline, err := client.cachedPipeline(svcConfig)
	if err != nil {
		return nil, err
	}

	resp, err := pipeline.Do(req)
	if err != nil {
		return nil, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}

	// unmarshal response
	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *GeneralClient) CreateOrUpdateThenPoll(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, body any, options RequestOptions) (interface{}, error) {
	// build request
	urlPath := fmt.Sprintf("https://%s", url)
	req, err := runtime.NewRequest(ctx, http.MethodPut, urlPath)
	if err != nil {
		return nil, err
	}

	reqQP := req.Raw().URL.Query()
	if apiVersion != "" {
		reqQP.Set(apiVersionHeaderKey, apiVersion)
	}

	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}

	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}

	err = runtime.MarshalAsJSON(req, body)
	if err != nil {
		return nil, err
	}

	// send request
	pipeline, err := client.cachedPipeline(svcConfig)
	if err != nil {
		return nil, err
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}

	// poll until done
	pt, err := runtime.NewPoller[interface{}](resp, pipeline, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}

	// unmarshal response
	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *GeneralClient) Action(ctx context.Context, url, apiVersion, method string, svcConfig cloud.ServiceConfiguration, body any, options RequestOptions) (interface{}, error) {
	// build request
	urlPath := fmt.Sprintf("https://%s", url)

	if method == "" {
		method = http.MethodGet
	}

	req, err := runtime.NewRequest(ctx, method, urlPath)
	if err != nil {
		return nil, err
	}

	reqQP := req.Raw().URL.Query()

	if apiVersion != "" {
		reqQP.Set(apiVersionHeaderKey, apiVersion)
	}
	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()

	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}

	if method != http.MethodGet && body != nil {
		err = runtime.MarshalAsJSON(req, body)
	}

	if err != nil {
		return nil, err
	}

	// send request
	pipeline, err := client.cachedPipeline(svcConfig)
	if err != nil {
		return nil, err
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}

	// poll until done
	pt, err := runtime.NewPoller[interface{}](resp, pipeline, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}

	// unmarshal response
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

func (client *GeneralClient) DeleteThenPoll(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, options RequestOptions) (interface{}, error) {
	// build request
	urlPath := fmt.Sprintf("https://%s", url)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, urlPath)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()

	if apiVersion != "" {
		reqQP.Set("api-version", apiVersion)
	}

	for key, value := range options.QueryParameters {
		reqQP.Set(key, value)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()

	req.Raw().Header.Set("Accept", "application/json")
	for key, value := range options.Headers {
		req.Raw().Header.Set(key, value)
	}

	// send request
	pipeline, err := client.cachedPipeline(svcConfig)
	if err != nil {
		return nil, err
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}

	// poll until done
	pt, err := runtime.NewPoller[interface{}](resp, pipeline, nil)
	if err == nil {
		resp, err := pt.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 10 * time.Second,
		})
		return resp, err
	}

	// unmarshal response
	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (retryclient *GeneralClientRetryableErrors) CreateOrUpdateThenPoll(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, body any, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	ctx = tflog.SetField(ctx, "request", "CreateOrUpdateThenPoll")
	ctx = retryclient.updateContext(ctx)
	tflog.Debug(ctx, "retryclient: Begin")
	i := 0
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.CreateOrUpdateThenPoll(ctx, url, apiVersion, svcConfig, body, options)
			if err != nil {
				if isGeneralRetryable(ctx, *retryclient, data, err) {
					tflog.Debug(ctx, "retryclient: Retry attempt", map[string]interface{}{
						"err":     err,
						"attempt": i,
					})
					i++
					return data, err
				}
				tflog.Debug(ctx, "retryclient: PermanentError", map[string]interface{}{
					"err":     err,
					"attempt": i,
				})
				tflog.Debug(ctx, "retryclient: Success", map[string]interface{}{
					"attempt": i,
				})
				return nil, &backoff.PermanentError{Err: err}
			}
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData(op, exbo)
}

func (retryclient *GeneralClientRetryableErrors) Get(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	ctx = tflog.SetField(ctx, "request", "Get")
	ctx = retryclient.updateContext(ctx)
	tflog.Debug(ctx, "retryclient: Begin")
	i := 0
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.Get(ctx, url, apiVersion, svcConfig, options)
			if err != nil {
				if isGeneralRetryable(ctx, *retryclient, data, err) {
					tflog.Debug(ctx, "retryclient: Retry attempt", map[string]interface{}{
						"err":     err,
						"attempt": i,
					})
					i++
					return data, err
				}
				tflog.Debug(ctx, "retryclient: PermanentError", map[string]interface{}{
					"err":     err,
					"attempt": i,
				})
				return nil, &backoff.PermanentError{Err: err}
			}
			tflog.Debug(ctx, "retryclient: Success", map[string]interface{}{
				"attempt": i,
			})
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData(op, exbo)
}

func (retryclient *GeneralClientRetryableErrors) DeleteThenPoll(ctx context.Context, url, apiVersion string, svcConfig cloud.ServiceConfiguration, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	ctx = tflog.SetField(ctx, "request", "DeleteThenPoll")
	ctx = retryclient.updateContext(ctx)
	tflog.Debug(ctx, "retryclient: Begin")
	i := 0
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.DeleteThenPoll(ctx, url, apiVersion, svcConfig, options)
			if err != nil {
				if isGeneralRetryable(ctx, *retryclient, data, err) {
					tflog.Debug(ctx, "retryclient: Retry attempt", map[string]interface{}{
						"err":     err,
						"attempt": i,
					})
					i++
					return data, err
				}
				tflog.Debug(ctx, "retryclient: PermanentError", map[string]interface{}{
					"err":     err,
					"attempt": i,
				})
				return nil, &backoff.PermanentError{Err: err}
			}
			tflog.Debug(ctx, "retryclient: Success", map[string]interface{}{
				"attempt": i,
			})
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData(op, exbo)
}

func (retryclient *GeneralClientRetryableErrors) Action(ctx context.Context, url, apiVersion, method string, svcConfig cloud.ServiceConfiguration, body any, options RequestOptions) (interface{}, error) {
	if retryclient.backoff == nil {
		return nil, errors.New("retry is not configured, please call WithRetry() first")
	}
	ctx = tflog.SetField(ctx, "request", "Action")
	ctx = retryclient.updateContext(ctx)
	tflog.Debug(ctx, "retryclient: Begin")
	i := 0
	op := backoff.OperationWithData[interface{}](
		func() (interface{}, error) {
			data, err := retryclient.client.Action(ctx, url, apiVersion, method, svcConfig, body, options)
			if err != nil {
				if isGeneralRetryable(ctx, *retryclient, data, err) {
					tflog.Debug(ctx, "retryclient: Retry attempt", map[string]interface{}{
						"err":     err,
						"attempt": i,
					})
					i++
					return data, err
				}
				tflog.Debug(ctx, "retryclient: PermanentError", map[string]interface{}{
					"err":     err,
					"attempt": i,
				})
				return nil, &backoff.PermanentError{Err: err}
			}
			tflog.Debug(ctx, "retryclient: Success", map[string]interface{}{
				"attempt": i,
			})
			return data, err
		})
	exbo := backoff.WithContext(retryclient.backoff, ctx)
	return backoff.RetryWithData(op, exbo)
}

func isGeneralRetryable(ctx context.Context, retryclient GeneralClientRetryableErrors, data interface{}, err error) bool {
	for _, e := range retryclient.errors {
		if e.MatchString(err.Error()) {
			tflog.Debug(ctx, "isDataPlaneRetryable: Error is retryable by regex", map[string]interface{}{
				"err":    err,
				"regexp": e.String(),
			})
			return true
		}
	}
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		if slices.Contains(retryclient.statusCodes, respErr.StatusCode) {
			tflog.Debug(ctx, "isDataPlaneRetryable: Error is retryable by status code", map[string]interface{}{
				"err":        err,
				"statusCode": respErr.StatusCode,
			})
			return true
		}
	}
	for i, f := range retryclient.dataCallbackFuncs {
		if f(data) {
			tflog.Debug(ctx, "isDataPlaneRetryable: Error is retryable by function callback", map[string]interface{}{
				"err":               err,
				"callback_func_idx": i,
			})
			return true
		}
	}
	return false
}

// ConfigureClientWithCustomRetry configures the client with a custom retry configuration if supplied.
// If the retry configuration is null or unknown, it will use the default retry configuration.
// If the supplied context has a deadline, it will use the deadline as the max elapsed time when a custom retry is provided.
func (client *GeneralClient) ConfigureClientWithCustomRetry(ctx context.Context, rtry retry.RetryValue, useReadAfterCreateValues bool) GeneralRequester {
	backOff, errRegExps, statusCodes := configureCustomRetry(ctx, rtry, useReadAfterCreateValues)
	return client.WithRetry(backOff, errRegExps, statusCodes, nil)
}
