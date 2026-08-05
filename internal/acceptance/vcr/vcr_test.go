package vcr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
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
