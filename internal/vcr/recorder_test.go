package vcr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

func TestModeFromEnv(t *testing.T) {
	cases := []struct {
		value string
		want  Mode
	}{
		{"", ModePassthrough},
		{"off", ModePassthrough},
		{"RECORD", ModeRecord},
		{" record", ModeRecord},
		{"replay", ModeReplay},
		{"true", ModeReplay},
		{"nonsense", ModePassthrough},
	}
	for _, tc := range cases {
		t.Setenv(EnvVar, tc.value)
		if got := ModeFromEnv(); got != tc.want {
			t.Errorf("ModeFromEnv(%q) = %q, want %q", tc.value, got, tc.want)
		}
	}
}

func TestModeHelpers(t *testing.T) {
	t.Setenv(EnvVar, "replay")
	if !Enabled() || !IsReplaying() || IsRecording() {
		t.Fatalf("replay: Enabled=%v IsReplaying=%v IsRecording=%v", Enabled(), IsReplaying(), IsRecording())
	}
	t.Setenv(EnvVar, "record")
	if !Enabled() || IsReplaying() || !IsRecording() {
		t.Fatalf("record: Enabled=%v IsReplaying=%v IsRecording=%v", Enabled(), IsReplaying(), IsRecording())
	}
	t.Setenv(EnvVar, "")
	if Enabled() || IsReplaying() || IsRecording() {
		t.Fatalf("passthrough: Enabled=%v IsReplaying=%v IsRecording=%v", Enabled(), IsReplaying(), IsRecording())
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

func TestFakeCredentialTokenClaims(t *testing.T) {
	tok, err := FakeCredential().GetToken(context.Background(), policy.TokenRequestOptions{})
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

func TestSubscriptionMatcher(t *testing.T) {
	const canonicalURL = "https://management.azure.com/subscriptions/" + CanonicalSubscriptionID + "/resourceGroups/rg?api-version=2023-07-01"
	i := cassette.Request{Method: http.MethodGet, URL: canonicalURL}

	// A request carrying a real subscription ID still matches a canonical cassette.
	realURL := "https://management.azure.com/subscriptions/aaaaaaaa-1111-2222-3333-444444444444/resourceGroups/rg?api-version=2023-07-01"
	if !subscriptionMatcher(httptest.NewRequest(http.MethodGet, realURL, nil), i) {
		t.Error("expected request with real subscription ID to match canonical cassette")
	}
	// Method mismatch does not match.
	if subscriptionMatcher(httptest.NewRequest(http.MethodPut, canonicalURL, nil), i) {
		t.Error("expected method mismatch to fail matching")
	}
	// Different path does not match.
	other := "https://management.azure.com/subscriptions/" + CanonicalSubscriptionID + "/resourceGroups/other?api-version=2023-07-01"
	if subscriptionMatcher(httptest.NewRequest(http.MethodGet, other, nil), i) {
		t.Error("expected different path to fail matching")
	}
}

func TestGetRecorderRequiresTestName(t *testing.T) {
	if _, err := GetRecorder("", ""); err == nil {
		t.Fatal("expected GetRecorder to reject an empty test name")
	}
}

func TestGetRecorderIsIdempotent(t *testing.T) {
	// Passthrough mode does not persist a cassette, so this stays hermetic.
	t.Setenv(EnvVar, "")
	const name = "TestGetRecorderIsIdempotent"
	t.Cleanup(func() {
		_ = StopRecorder(name)
		_ = os.Remove(filepath.Join(testDataPath, name+".yaml"))
	})

	first, err := GetRecorder(name, "")
	if err != nil {
		t.Fatalf("GetRecorder: %v", err)
	}
	second, err := GetRecorder(name, "")
	if err != nil {
		t.Fatalf("GetRecorder (second): %v", err)
	}
	if first != second {
		t.Fatal("GetRecorder should return the same recorder for the same test name")
	}
	if err := StopRecorder(name); err != nil {
		t.Fatalf("StopRecorder: %v", err)
	}
	// After stopping, a subsequent GetRecorder yields a fresh recorder.
	third, err := GetRecorder(name, "")
	if err != nil {
		t.Fatalf("GetRecorder (third): %v", err)
	}
	if third == first {
		t.Fatal("GetRecorder should create a new recorder after StopRecorder")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRecordThenReplayRoundTrip(t *testing.T) {
	name := filepath.Join(t.TempDir(), "cassette")
	const url = "https://management.azure.com/subscriptions/" + CanonicalSubscriptionID + "/resource"

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
		recorder.WithMatcher(subscriptionMatcher),
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
		recorder.WithMatcher(subscriptionMatcher),
	)
	if err != nil {
		t.Fatalf("creating replay recorder: %v", err)
	}
	defer func() { _ = replay.Stop() }()

	// A request with a real subscription ID replays against the canonical cassette.
	realURL := "https://management.azure.com/subscriptions/aaaaaaaa-1111-2222-3333-444444444444/resource"
	resp2, err := replay.GetDefaultClient().Get(realURL)
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
