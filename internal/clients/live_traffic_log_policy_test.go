package clients

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestLiveTrafficLogPolicyLogging(t *testing.T) {
	const requestBody = `{"secret":"request-secret"}`
	const responseBody = `{"secret":"response-secret"}`
	for _, tc := range []struct {
		name    string
		env     string
		logBody bool
	}{
		{name: "default"},
		{name: "disabled", env: "false"},
		{name: "invalid", env: "invalid"},
		{name: "enabled", env: "true", logBody: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("LOG_SENSITIVE_DATA", tc.env)
			var logOutput bytes.Buffer
			originalOutput, originalFlags := log.Writer(), log.Flags()
			log.SetOutput(&logOutput)
			log.SetFlags(0)
			t.Cleanup(func() {
				log.SetOutput(originalOutput)
				log.SetFlags(originalFlags)
			})

			pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{
				PerRetry: []policy.Policy{NewLiveTrafficLogPolicy()},
			}, &policy.ClientOptions{
				Telemetry: policy.TelemetryOptions{Disabled: true},
				Transport: fakeTransporter{response: &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(responseBody)),
				}},
			})
			req, err := runtime.NewRequest(context.Background(), http.MethodPut, "https://example.com")
			if err != nil {
				t.Fatal(err)
			}
			if err := runtime.MarshalAsJSON(req, json.RawMessage(requestBody)); err != nil {
				t.Fatal(err)
			}
			req.Raw().Header.Set("Authorization", "Bearer auth-secret")
			req.Raw().Header.Set("x-ms-authorization-auxiliary", "Bearer auxiliary-secret")

			resp, err := pl.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(resp.Request.Body)
			if err != nil || string(body) != requestBody {
				t.Fatalf("request body was not preserved: body=%q, error=%v", body, err)
			}
			body, err = runtime.Payload(resp)
			if err != nil || string(body) != responseBody {
				t.Fatalf("response body was not preserved: body=%q, error=%v", body, err)
			}

			var got traffic
			payload := strings.TrimPrefix(logOutput.String(), "[DEBUG] Live traffic: ")
			if err := json.Unmarshal([]byte(payload), &got); err != nil {
				t.Fatal(err)
			}
			wantRequest, wantResponse := redactedValue, redactedValue
			if tc.logBody {
				wantRequest, wantResponse = requestBody, responseBody
			}
			if got.LiveRequest.Body != wantRequest || got.LiveResponse.Body != wantResponse {
				t.Fatalf("unexpected logged bodies: request=%q response=%q", got.LiveRequest.Body, got.LiveResponse.Body)
			}
			if got.LiveRequest.Headers["Authorization"] != redactedValue || got.LiveRequest.Headers["X-Ms-Authorization-Auxiliary"] != redactedValue {
				t.Fatal("authentication headers were not redacted")
			}
		})
	}
}

type fakeTransporter struct {
	response *http.Response
}

func (f fakeTransporter) Do(req *http.Request) (*http.Response, error) {
	f.response.Request = req
	return f.response, nil
}
