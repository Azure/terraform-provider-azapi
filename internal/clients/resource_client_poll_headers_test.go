package clients

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type pollHeadersTransport struct {
	host         string
	resourcePath string
	asyncOpPath  string
	headerName   string
	headerValue  string

	mu       sync.Mutex
	requests []string
}

func (t *pollHeadersTransport) Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get(t.headerName) != t.headerValue {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"error":{"code":"MissingRequiredHeader"}}`)),
			Request:    req,
		}, nil
	}

	t.mu.Lock()
	t.requests = append(t.requests, fmt.Sprintf("%s %s", req.Method, req.URL.Path))
	t.mu.Unlock()

	switch {
	case req.Method == http.MethodPut && req.URL.Path == t.resourcePath:
		resp := &http.Response{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"properties":{"provisioningState":"Creating"}}`)),
			Request:    req,
		}
		resp.Header.Set("Azure-AsyncOperation", t.host+t.asyncOpPath)
		return resp, nil

	case req.Method == http.MethodGet && req.URL.Path == t.asyncOpPath:
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"status":"Succeeded"}`)),
			Request:    req,
		}, nil

	case req.Method == http.MethodGet && req.URL.Path == t.resourcePath:
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"properties":{"provisioningState":"Succeeded"}}`)),
			Request:    req,
		}, nil
	}

	return nil, fmt.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
}

func (t *pollHeadersTransport) requestLog() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]string(nil), t.requests...)
}

func TestResourceClientCreateOrUpdate_PropagatesHeadersToPollRequests(t *testing.T) {
	const (
		host         = "https://management.azure.com"
		resourcePath = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/virtualMachineScaleSets/example"
		asyncOpPath  = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/locations/westeurope/operations/11111111-1111-1111-1111-111111111111"
		headerName   = "X-Ms-Target-Resource-Id"
		headerValue  = "/subscriptions/target-subscription/resourceGroups/target-rg"
	)

	transport := &pollHeadersTransport{
		host:         host,
		resourcePath: resourcePath,
		asyncOpPath:  asyncOpPath,
		headerName:   headerName,
		headerValue:  headerValue,
	}
	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{Disabled: true},
		Transport: transport,
	})
	client := &ResourceClient{host: host, pl: pl}

	_, err := client.CreateOrUpdate(context.Background(), resourcePath, "2023-09-01", map[string]interface{}{}, RequestOptions{
		Headers: map[string]string{headerName: headerValue},
	})
	if err != nil {
		t.Fatalf("CreateOrUpdate should send custom headers on every LRO request: %v", err)
	}

	want := []string{
		http.MethodPut + " " + resourcePath,
		http.MethodGet + " " + asyncOpPath,
		http.MethodGet + " " + resourcePath,
	}
	got := transport.requestLog()
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected request sequence:\n got: %v\nwant: %v", got, want)
	}
}
