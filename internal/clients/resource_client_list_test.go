package clients

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type listTransport struct {
	responses map[string]*http.Response
	errByPath map[string]error
	requests  []*http.Request
}

func (t *listTransport) Do(req *http.Request) (*http.Response, error) {
	t.requests = append(t.requests, req)
	key := req.URL.Path
	if req.URL.RawQuery != "" {
		key = fmt.Sprintf("%s?%s", key, req.URL.RawQuery)
	}

	if err, ok := t.errByPath[key]; ok {
		return nil, err
	}

	if resp, ok := t.responses[key]; ok {
		clone := *resp
		clone.Request = req
		return &clone, nil
	}

	return nil, fmt.Errorf("unexpected request path %q", key)
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestResourceClientListAggregatesPages(t *testing.T) {
	transport := &listTransport{
		responses: map[string]*http.Response{
			"/resources?api-version=2021-01-01": jsonResponse(http.StatusOK, `{"value":[{"id":"/r1"}],"nextLink":"https://example.com/resources?page=2"}`),
			"/resources?page=2":                 jsonResponse(http.StatusOK, `{"value":[{"id":"/r2"}]}`),
		},
		errByPath: map[string]error{},
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{Disabled: true},
		Transport: transport,
	})

	client := &ResourceClient{
		host: "https://example.com",
		pl:   pl,
	}

	output, err := client.List(context.Background(), "/resources", "2021-01-01", RequestOptions{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	resultMap, ok := output.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map output, got %T", output)
	}

	values, ok := resultMap["value"].([]interface{})
	if !ok {
		t.Fatalf("expected value to be array, got %T", resultMap["value"])
	}

	if len(values) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(values))
	}
}

func TestResourceClientListReturnsErrorOnPagingFailure(t *testing.T) {
	transport := &listTransport{
		responses: map[string]*http.Response{
			"/resources?api-version=2021-01-01": jsonResponse(http.StatusOK, `{"value":[{"id":"/r1"}],"nextLink":"https://example.com/resources?page=2"}`),
		},
		errByPath: map[string]error{
			"/resources?page=2": errors.New("simulated high-latency paging timeout"),
		},
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{Disabled: true},
		Transport: transport,
	})

	client := &ResourceClient{
		host: "https://example.com",
		pl:   pl,
	}

	output, err := client.List(context.Background(), "/resources", "2021-01-01", RequestOptions{})
	if err == nil {
		t.Fatalf("expected error when paging fails")
	}
	if output != nil {
		t.Fatalf("expected nil output when paging fails, got %#v", output)
	}
}
