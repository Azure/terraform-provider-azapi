package vcr

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

func TestSanitizeInteraction(t *testing.T) {
	const (
		realSub    = "aaaaaaaa-1111-2222-3333-444444444444"
		realTenant = "bbbbbbbb-5555-6666-7777-888888888888"
	)

	i := &cassette.Interaction{
		Request: cassette.Request{
			Method: http.MethodGet,
			URL:    "https://management.azure.com/subscriptions/" + realSub + "/x",
			Body:   `{"tenantId":"` + realTenant + `"}`,
			Headers: http.Header{
				"Authorization":                []string{"******"},
				"X-Ms-Authorization-Auxiliary": []string{"******"},
				"X-Custom":                     []string{"sub " + realSub},
			},
		},
		Response: cassette.Response{
			Code: 200,
			Body: `{"id":"/subscriptions/` + realSub + `/rg"}`,
			Headers: http.Header{
				"Set-Cookie": []string{"session=abc"},
			},
		},
	}

	if err := sanitizeInteraction(realSub, realTenant)(i); err != nil {
		t.Fatalf("sanitizeInteraction: %v", err)
	}

	if got := i.Request.Headers.Get("Authorization"); got != "" {
		t.Errorf("Authorization header should be removed, got %q", got)
	}
	if got := i.Request.Headers.Get("X-Ms-Authorization-Auxiliary"); got != "" {
		t.Errorf("aux auth header should be removed, got %q", got)
	}
	if got := i.Response.Headers.Get("Set-Cookie"); got != "" {
		t.Errorf("Set-Cookie should be removed, got %q", got)
	}
	for _, blob := range []string{i.Request.URL, i.Request.Body, i.Response.Body, i.Request.Headers.Get("X-Custom")} {
		if strings.Contains(blob, realSub) {
			t.Errorf("real subscription id not redacted in %q", blob)
		}
		if strings.Contains(blob, realTenant) {
			t.Errorf("real tenant id not redacted in %q", blob)
		}
	}
	if !strings.Contains(i.Request.URL, CanonicalSubscriptionID) {
		t.Errorf("canonical subscription id missing from URL: %q", i.Request.URL)
	}
	if !strings.Contains(i.Request.Body, CanonicalTenantID) {
		t.Errorf("canonical tenant id missing from body: %q", i.Request.Body)
	}
}

func TestSanitizeInteractionRemovesTracingHeaders(t *testing.T) {
	i := &cassette.Interaction{
		Request: cassette.Request{
			Headers: http.Header{
				"X-Ms-Client-Request-Id": {"client-request-id"},
				"Traceparent":            {"trace-id"},
			},
		},
		Response: cassette.Response{
			Headers: http.Header{
				"X-Ms-Correlation-Request-Id": {"correlation-request-id"},
				"X-Ms-Request-Id":             {"request-id"},
				"X-Ms-Routing-Request-Id":     {"routing-request-id"},
				"X-Ms-Arm-Service-Request-Id": {"arm-service-request-id"},
				"X-Ms-Operation-Identifier":   {"operation-identifier"},
				"X-Msedge-Ref":                {"edge-ref"},
				"Azure-Asyncoperation":        {"https://example.com/operation"},
			},
		},
	}

	if err := sanitizeInteraction("", "")(i); err != nil {
		t.Fatalf("sanitizeInteraction: %v", err)
	}

	for _, header := range []string{
		"X-Ms-Client-Request-Id",
		"Traceparent",
		"X-Ms-Correlation-Request-Id",
		"X-Ms-Request-Id",
		"X-Ms-Routing-Request-Id",
		"X-Ms-Arm-Service-Request-Id",
		"X-Ms-Operation-Identifier",
		"X-Msedge-Ref",
	} {
		if got := i.Request.Headers.Get(header); got != "" {
			t.Errorf("request header %s should be removed, got %q", header, got)
		}
		if got := i.Response.Headers.Get(header); got != "" {
			t.Errorf("response header %s should be removed, got %q", header, got)
		}
	}
	if got := i.Response.Headers.Get("Azure-Asyncoperation"); got == "" {
		t.Error("Azure-Asyncoperation header should be preserved")
	}
}

func TestSanitizeInteractionRedactsTransientValues(t *testing.T) {
	const (
		pid          = "222c6c49-1b0a-5959-a213-6608f9eb8820"
		etag         = "700020e8-b21e-4d09-8d86-887efa7abb5f"
		resourceGUID = "35185fcd-a674-41c5-bf9d-15b616785c89"
	)
	i := &cassette.Interaction{
		Request: cassette.Request{
			URL:        "https://example.com/operation?api-version=1&t=time&c=certificate&h=hash&s=signature",
			RequestURI: "/operation?api-version=1&t=time&c=certificate&h=hash&s=signature",
			Form: url.Values{
				"api-version": {"1"},
				"c":           {"certificate"},
				"h":           {"hash"},
				"s":           {"signature"},
				"t":           {"time"},
			},
			Headers: http.Header{
				"User-Agent": {"terraform-provider-azapi/acc pid-" + pid},
			},
		},
		Response: cassette.Response{
			ContentLength: 1,
			Body:          `{"etag":"W/\"` + etag + `\"","resourceGuid":"` + resourceGUID + `","resourceGuide":"` + resourceGUID + `"}`,
			Headers: http.Header{
				"Azure-Asyncoperation": {"https://example.com/operation?api-version=1&t=time&c=certificate&h=hash&s=signature"},
				"Content-Length":       {"1"},
			},
		},
	}

	if err := sanitizeInteraction("", "")(i); err != nil {
		t.Fatalf("sanitizeInteraction: %v", err)
	}

	for _, value := range []string{
		i.Request.URL,
		i.Request.RequestURI,
		i.Request.Headers.Get("User-Agent"),
		i.Response.Body,
		i.Response.Headers.Get("Azure-Asyncoperation"),
	} {
		for _, sensitive := range []string{pid, etag, resourceGUID, "certificate", "hash", "signature", "time"} {
			if strings.Contains(value, sensitive) {
				t.Errorf("sensitive value %q not redacted in %q", sensitive, value)
			}
		}
	}
	if !strings.Contains(i.Request.URL, "api-version=1") {
		t.Errorf("non-sensitive query parameter was changed in %q", i.Request.URL)
	}
	for _, key := range []string{"c", "h", "s", "t"} {
		if got := i.Request.Form.Get(key); got != redactedValue {
			t.Errorf("form value %s = %q, want %q", key, got, redactedValue)
		}
	}
	if got := i.Request.Form.Get("api-version"); got != "1" {
		t.Errorf("api-version form value = %q, want 1", got)
	}
	if got, want := i.Response.ContentLength, int64(len(i.Response.Body)); got != want {
		t.Errorf("response content length = %d, want %d", got, want)
	}
	if got, want := i.Response.Headers.Get("Content-Length"), strconv.Itoa(len(i.Response.Body)); got != want {
		t.Errorf("Content-Length header = %q, want %q", got, want)
	}
}
