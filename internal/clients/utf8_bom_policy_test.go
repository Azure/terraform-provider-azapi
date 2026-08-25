package clients

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func TestUTF8BOMPolicyJSONResponse(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{
			name: "removes BOM",
			body: "\xef\xbb\xbf{\"value\":\"expected\"}",
		},
		{
			name: "preserves body without BOM",
			body: "{\"value\":\"expected\"}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := runUTF8BOMPolicy(t, "application/json; charset=utf-8", testCase.body)

			var got map[string]string
			if err := runtime.UnmarshalAsJSON(resp, &got); err != nil {
				t.Fatalf("unmarshalling response: %+v", err)
			}
			if got["value"] != "expected" {
				t.Fatalf("expected decoded value, got %#v", got)
			}
			if resp.ContentLength != int64(len(`{"value":"expected"}`)) {
				t.Fatalf("expected adjusted content length, got %d", resp.ContentLength)
			}
		})
	}
}

func TestUTF8BOMPolicyPreservesNonJSONResponse(t *testing.T) {
	const body = "\xef\xbb\xbfplain text"
	resp := runUTF8BOMPolicy(t, "text/plain", body)

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response: %+v", err)
	}
	if string(got) != body {
		t.Fatalf("expected body %q, got %q", body, got)
	}
	if resp.ContentLength != int64(len(body)) {
		t.Fatalf("expected unchanged content length, got %d", resp.ContentLength)
	}
}

func runUTF8BOMPolicy(t *testing.T, contentType string, body string) *http.Response {
	t.Helper()

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			newUTF8BOMPolicy(),
		},
	}, &policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{Disabled: true},
		Transport: fakeTransporter{
			response: &http.Response{
				StatusCode:    http.StatusOK,
				Header:        http.Header{"Content-Type": []string{contentType}},
				Body:          io.NopCloser(strings.NewReader(body)),
				ContentLength: int64(len(body)),
			},
		},
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.com")
	if err != nil {
		t.Fatalf("creating request: %+v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("executing request: %+v", err)
	}
	return resp
}
