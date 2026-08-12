package vcr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// EnvVar controls the VCR framework. It is deliberately named to match the
// sibling terraform-provider-azurerm implementation so the two providers share
// a mental model:
//
//	TC_TEST_VIA_VCR=record        record live traffic to a cassette
//	TC_TEST_VIA_VCR=replay|true   replay from a cassette, no network or credentials
//	unset / anything else         passthrough (VCR disabled)
const EnvVar = "TC_TEST_VIA_VCR"

type Mode string

const (
	ModePassthrough Mode = "passthrough"
	ModeRecord      Mode = "record"
	ModeReplay      Mode = "replay"
)

// ModeFromEnv reports the VCR mode requested via EnvVar.
func ModeFromEnv() Mode {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvVar))) {
	case "record":
		return ModeRecord
	case "replay", "true":
		return ModeReplay
	default:
		return ModePassthrough
	}
}

// Enabled reports whether VCR is recording or replaying (i.e. not passthrough).
func Enabled() bool { return ModeFromEnv() != ModePassthrough }

// IsReplaying reports whether VCR is replaying from cassettes.
func IsReplaying() bool { return ModeFromEnv() == ModeReplay }

// IsRecording reports whether VCR is recording live traffic.
func IsRecording() bool { return ModeFromEnv() == ModeRecord }

// Canonical placeholders that real identifiers are rewritten to before a
// cassette is saved. They are also the values the fake credential mints during
// replay so recorded and replayed requests agree byte-for-byte.
const (
	CanonicalSubscriptionID = "00000000-0000-0000-0000-000000000000"
	CanonicalTenantID       = "11111111-1111-1111-1111-111111111111"
	CanonicalObjectID       = "22222222-2222-2222-2222-222222222222"
	CanonicalClientID       = "33333333-3333-3333-3333-333333333333"
	CanonicalClientSecret   = "fake-client-secret"
)

// Fixed, name-derived locations used by VCR tests so requests are deterministic.
const (
	LocationPrimary   = "eastus"
	LocationSecondary = "westeurope"
	LocationTernary   = "eastasia"
)

// testDataPath is the directory (relative to the test package) that cassettes
// are read from and written to.
const testDataPath = "testdata/cassettes"

// testRecorder pairs a go-vcr recorder with the sanitizer that scrubs its
// interactions, so runtime-registered secrets (see RegisterSecret) reach the
// BeforeSave hook.
type testRecorder struct {
	rec *recorder.Recorder
	san *sanitizer
}

var (
	recorders = make(map[string]*testRecorder)
	mu        sync.Mutex
)

// GetRecorder returns the shared recorder for a given test name, initialising it
// if necessary. It redacts sensitive information (subscription/tenant IDs,
// credentials) via the sanitizer's BeforeSave hook and tailors the matcher to
// azapi requests. subscriptionID is the subscription the client under test uses;
// it is redacted to CanonicalSubscriptionID in addition to the values sourced
// from the environment.
func GetRecorder(testName string, subscriptionID string) (*recorder.Recorder, error) {
	if testName == "" {
		return nil, errors.New("testName must be provided to retrieve a recorder")
	}

	mu.Lock()
	defer mu.Unlock()

	if tr, exists := recorders[testName]; exists {
		return tr.rec, nil
	}

	// default to passthrough, just in case something unexpected is set
	mode := recorder.ModePassthrough
	switch ModeFromEnv() {
	case ModeRecord:
		mode = recorder.ModeRecordOnly
	case ModeReplay:
		mode = recorder.ModeReplayOnly
	}

	// The sanitizer redacts known sensitive values and then audits every
	// interaction, failing the recording (fail-closed) if a secret survives.
	san := newSanitizer(os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_TENANT_ID"))
	if subscriptionID != "" {
		san.AddReplacement(subscriptionID, CanonicalSubscriptionID)
	}
	// Canonicalize non-secret identifiers and redact the credential material the
	// harness itself supplies, so none of it can be written to a cassette.
	san.AddReplacement(os.Getenv("ARM_CLIENT_ID"), CanonicalClientID)
	san.AddReplacement(os.Getenv("ARM_READER_CLIENT_ID"), CanonicalClientID)
	san.AddSecret(os.Getenv("ARM_CLIENT_SECRET"))
	san.AddSecret(os.Getenv("ARM_READER_CLIENT_SECRET"))
	san.AddSecret(os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"))
	san.AddSecret(os.Getenv("ARM_OIDC_TOKEN"))

	cassettePath := filepath.Join(testDataPath, cassetteName(testName))
	rec, err := recorder.New(cassettePath,
		recorder.WithMode(mode),
		recorder.WithSkipRequestLatency(true),
		recorder.WithMatcher(subscriptionMatcher),
		// recorder.WithFS(&GzipFS{}),
		recorder.WithHook(san.BeforeSave, recorder.BeforeSaveHook),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create recorder for %s: %w", testName, err)
	}

	recorders[testName] = &testRecorder{rec: rec, san: san}
	return rec, nil
}

// StopRecorder stops and removes the recorder from the map, saving it to disk.
func StopRecorder(testName string) error {
	mu.Lock()
	defer mu.Unlock()

	if tr, exists := recorders[testName]; exists {
		err := tr.rec.Stop()
		delete(recorders, testName)
		if err != nil {
			return fmt.Errorf("failed to stop recorder for %s: %w", testName, err)
		}
	}
	return nil
}

// RegisterSecret registers a sensitive value with the recorder for the test
// running on the current goroutine so it is scrubbed from the cassette before it
// is saved. It is safe to call unconditionally: when VCR is disabled (or no
// recorder is active for this goroutine) it does nothing. Use it from tests that
// read secret values back from Azure, e.g. vcr.RegisterSecret(storageAccountKey).
func RegisterSecret(v string) {
	if v == "" {
		return
	}
	testName := CurrentTestName()
	if testName == "" {
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if tr, exists := recorders[testName]; exists {
		tr.san.AddSecret(v)
	}
}

// FakeCredential returns a credential that mints a syntactically valid token
// with canonical claims without contacting Azure. It is used during replay so
// the go-azure-sdk pipeline never performs a real token exchange.
func FakeCredential() azcore.TokenCredential { return fakeCredential{} }

// cassetteName maps a test name to its on-disk cassette base name (go-vcr
// appends the .yaml extension). Slashes from subtests are flattened.
func cassetteName(testName string) string {
	return strings.ReplaceAll(testName, "/", "_")
}

// subscriptionURLPattern matches an Azure subscription ID inside a request URL.
var subscriptionURLPattern = regexp.MustCompile(`(?i)(/subscriptions/)([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

// normalizeSubscriptionIDs rewrites any subscription ID in s to the canonical
// placeholder. Cassettes are saved with canonical IDs, so normalising both the
// incoming request and the recorded interaction before comparison lets a request
// carrying a real subscription ID (e.g. during an ImportStep) match a redacted
// cassette.
func normalizeSubscriptionIDs(s string) string {
	return subscriptionURLPattern.ReplaceAllString(s, "${1}"+CanonicalSubscriptionID)
}

// subscriptionMatcher matches on method and URL after normalising subscription
// IDs on both sides.
func subscriptionMatcher(r *http.Request, i cassette.Request) bool {
	if r.Method != i.Method {
		return false
	}
	return normalizeSubscriptionIDs(r.URL.String()) == normalizeSubscriptionIDs(i.URL)
}

// Transport adapts a go-vcr recorder to the azcore policy.Transporter interface
// expected by the client options. The recorder's default client implements Do.
func Transport(rec *recorder.Recorder) policy.Transporter {
	if rec == nil {
		return nil
	}
	return rec.GetDefaultClient()
}

type fakeCredential struct{}

func (fakeCredential) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     fakeToken(),
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

func fakeToken() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	claims, _ := json.Marshal(map[string]any{
		"aud": "https://management.azure.com",
		"iss": "https://sts.windows.net/" + CanonicalTenantID + "/",
		"oid": CanonicalObjectID,
		"sub": CanonicalObjectID,
		"tid": CanonicalTenantID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Add(-time.Minute).Unix(),
	})
	payload := base64.RawURLEncoding.EncodeToString(claims)
	return header + "." + payload + "."
}

const charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"

func deterministicRand(name string) *rand.Rand {
	h := fnv.New64a()
	_, _ = h.Write([]byte(name))
	return rand.New(rand.NewSource(int64(h.Sum64()))) //nolint:gosec // determinism, not security
}

// RandInt returns a deterministic 18-digit integer derived from name, mirroring
// the shape of the live RandTimeInt while being stable across runs.
func RandInt(name string) int {
	r := deterministicRand(name + ":int")
	return int(r.Int63n(900000000000000000) + 100000000000000000)
}

// RandString returns a deterministic random string of the given length derived
// from name.
func RandString(name string, length int) string {
	r := deterministicRand(name + ":str")
	b := make([]byte, length)
	for i := range b {
		b[i] = charSetAlphaNum[r.Intn(len(charSetAlphaNum))]
	}
	return string(b)
}
