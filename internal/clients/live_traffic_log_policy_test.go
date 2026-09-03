package clients

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func TestLiveTrafficLogPolicyRedactsRequestBody(t *testing.T) {
	p := &liveTrafficLogPolicy{}
	req, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://example.com")
	if err != nil {
		t.Fatalf("failed to create request: %+v", err)
	}
	if err := runtime.MarshalAsJSON(req, map[string]string{"secret": "request-secret"}); err != nil {
		t.Fatalf("failed to marshal request body: %+v", err)
	}

	got := p.requestBodyString(req)
	if got != redactedValue {
		t.Fatalf("expected redacted request body, got %q", got)
	}
	if strings.Contains(got, "request-secret") {
		t.Fatalf("request body included secret")
	}

	if err := req.RewindBody(); err != nil {
		t.Fatalf("failed to rewind request body: %+v", err)
	}
	body, err := io.ReadAll(req.Raw().Body)
	if err != nil {
		t.Fatalf("failed to read request body: %+v", err)
	}
	if !strings.Contains(string(body), "request-secret") {
		t.Fatalf("request body was not preserved")
	}
}

func TestLiveTrafficLogPolicyRedactsResponseBody(t *testing.T) {
	p := &liveTrafficLogPolicy{}
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(`{"secret":"response-secret"}`)),
	}

	got := p.responseBodyString(resp)
	if got != redactedValue {
		t.Fatalf("expected redacted response body, got %q", got)
	}
	if strings.Contains(got, "response-secret") {
		t.Fatalf("response body included secret")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %+v", err)
	}
	if !strings.Contains(string(body), "response-secret") {
		t.Fatalf("response body was not preserved")
	}
}

func TestLiveTrafficLogPolicyEmptyBodies(t *testing.T) {
	p := &liveTrafficLogPolicy{}
	req, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://example.com")
	if err != nil {
		t.Fatalf("failed to create request: %+v", err)
	}

	if got := p.requestBodyString(req); got != "" {
		t.Fatalf("expected empty request body, got %q", got)
	}
	if got := p.responseBodyString(&http.Response{}); got != "" {
		t.Fatalf("expected empty response body, got %q", got)
	}
}

func TestLiveTrafficLogPolicyDoesNotLogBodies(t *testing.T) {
	var logOutput bytes.Buffer
	originalOutput := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&logOutput)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			NewLiveTrafficLogPolicy(),
		},
	}, &policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{Disabled: true},
		Transport: fakeTransporter{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"secret":"response-secret"}`)),
			},
		},
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://example.com")
	if err != nil {
		t.Fatalf("failed to create request: %+v", err)
	}
	if err := runtime.MarshalAsJSON(req, map[string]string{"secret": "request-secret"}); err != nil {
		t.Fatalf("failed to marshal request body: %+v", err)
	}

	if _, err := pl.Do(req); err != nil {
		t.Fatalf("unexpected pipeline error: %+v", err)
	}

	got := logOutput.String()
	if !strings.Contains(got, redactedValue) {
		t.Fatalf("expected log output to contain redaction marker, got %q", got)
	}
	if strings.Contains(got, "request-secret") {
		t.Fatalf("log output included request secret: %s", got)
	}
	if strings.Contains(got, "response-secret") {
		t.Fatalf("log output included response secret: %s", got)
	}
}

type fakeTransporter struct {
	response *http.Response
}

func (f fakeTransporter) Do(req *http.Request) (*http.Response, error) {
	f.response.Request = req
	return f.response, nil
}
