package vcr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

func TestModeFromEnv(t *testing.T) {
	cases := []struct {
		value string
		want  Mode
	}{
		{"", ModeOff},
		{"off", ModeOff},
		{"RECORD", ModeRecord},
		{" record", ModeRecord},
		{"replay", ModeReplay},
		{"nonsense", ModeOff},
	}
	for _, tc := range cases {
		t.Setenv(EnvVarMode, tc.value)
		if got := ModeFromEnv(); got != tc.want {
			t.Errorf("ModeFromEnv(%q) = %q, want %q", tc.value, got, tc.want)
		}
	}
}

func TestDeterministicRand(t *testing.T) {
	if got, want := RandInt("a"), RandInt("a"); got != want {
		t.Fatal("RandInt is not deterministic for the same name")
	}
	if got, want := RandString("a", 5), RandString("a", 5); got != want {
		t.Fatal("RandString is not deterministic for the same name")
	}
	if a, b := RandInt("a"), RandInt("b"); a == b {
		t.Fatal("RandInt should differ for different names")
	}
	if got := RandString("a", 8); len(got) != 8 {
		t.Fatalf("RandString length = %d, want 8", len(got))
	}
	// RandInt mirrors RandTimeInt: an 18-digit positive integer.
	if n := RandInt("something"); n < 100000000000000000 || n > 999999999999999999 {
		t.Fatalf("RandInt out of expected 18-digit range: %d", n)
	}
}

func TestFakeTokenClaims(t *testing.T) {
	tok, err := fakeCredential{}.GetToken(context.Background(), policy.TokenRequestOptions{})
	if err != nil {
		t.Fatalf("GetToken: %v", err)
	}
	parts := strings.Split(tok.Token, ".")
	if len(parts) != 3 {
		t.Fatalf("token should have 3 parts, got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decoding token payload: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("unmarshalling claims: %v", err)
	}
	if claims["oid"] != CanonicalObjectID {
		t.Errorf("oid = %v, want %v", claims["oid"], CanonicalObjectID)
	}
	if claims["tid"] != CanonicalTenantID {
		t.Errorf("tid = %v, want %v", claims["tid"], CanonicalTenantID)
	}
}

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
				"Authorization":                []string{"Bearer secret-token"},
				"X-Ms-Authorization-Auxiliary": []string{"Bearer aux-token"},
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRecordThenReplayRoundTrip(t *testing.T) {
	name := filepath.Join(t.TempDir(), "cassette")
	const url = "https://management.azure.com/subscriptions/x/resource"

	stub := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Request:    req,
		}, nil
	})

	rec, err := recorder.New(name,
		recorder.WithMode(recorder.ModeRecordOnly),
		recorder.WithRealTransport(stub),
		recorder.WithMatcher(methodURLMatcher),
	)
	if err != nil {
		t.Fatalf("creating record recorder: %v", err)
	}
	resp, err := rec.GetDefaultClient().Get(url)
	if err != nil {
		t.Fatalf("recording request: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if string(body) != `{"ok":true}` {
		t.Fatalf("recorded body = %q", body)
	}
	if err := rec.Stop(); err != nil {
		t.Fatalf("stopping record recorder: %v", err)
	}

	replay, err := recorder.New(name,
		recorder.WithMode(recorder.ModeReplayOnly),
		recorder.WithMatcher(methodURLMatcher),
	)
	if err != nil {
		t.Fatalf("creating replay recorder: %v", err)
	}
	defer func() { _ = replay.Stop() }()

	resp2, err := replay.GetDefaultClient().Get(url)
	if err != nil {
		t.Fatalf("replaying request: %v", err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	_ = resp2.Body.Close()
	if string(body2) != `{"ok":true}` {
		t.Fatalf("replayed body = %q, want the recorded body (no network)", body2)
	}

	if _, err := replay.GetDefaultClient().Get(url + "/does-not-exist"); err == nil {
		t.Fatal("expected an error replaying an unmatched request")
	}
}

func TestConfigureOptionReplayInjection(t *testing.T) {
	rec, err := recorder.New(filepath.Join(t.TempDir(), "cassette"),
		recorder.WithMode(recorder.ModeRecordOnly),
		recorder.WithRealTransport(roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, nil })),
	)
	if err != nil {
		t.Fatalf("creating recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	h := &Recorder{t: t, mode: ModeReplay, rec: rec}

	if !h.Enabled() {
		t.Fatal("harness should be enabled in replay mode")
	}
	if _, ok := h.Credential().(fakeCredential); !ok {
		t.Fatalf("Credential() = %T, want fakeCredential", h.Credential())
	}
	if h.SubscriptionID() != CanonicalSubscriptionID {
		t.Errorf("SubscriptionID() = %q, want canonical", h.SubscriptionID())
	}
	if h.TenantID() != CanonicalTenantID {
		t.Errorf("TenantID() = %q, want canonical", h.TenantID())
	}

	var opt clients.Option
	h.ConfigureOption(&opt)
	if opt.Transport == nil {
		t.Error("ConfigureOption should set a transport")
	}
	if _, ok := opt.Cred.(fakeCredential); !ok {
		t.Errorf("ConfigureOption Cred = %T, want fakeCredential", opt.Cred)
	}
	if opt.SubscriptionId != CanonicalSubscriptionID {
		t.Errorf("ConfigureOption SubscriptionId = %q, want canonical", opt.SubscriptionId)
	}
	if opt.TenantId != CanonicalTenantID {
		t.Errorf("ConfigureOption TenantId = %q, want canonical", opt.TenantId)
	}
}
