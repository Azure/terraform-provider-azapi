package vcr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"hash/fnv"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

type Mode string

const (
	ModeOff    Mode = "off"
	ModeRecord Mode = "record"
	ModeReplay Mode = "replay"
)

const EnvVarMode = "AZAPI_VCR_MODE"

const (
	CanonicalSubscriptionID = "00000000-0000-0000-0000-000000000000"
	CanonicalTenantID       = "11111111-1111-1111-1111-111111111111"
	CanonicalObjectID       = "22222222-2222-2222-2222-222222222222"
	CanonicalClientID       = "33333333-3333-3333-3333-333333333333"
	CanonicalClientSecret   = "fake-client-secret"
)

const (
	LocationPrimary   = "eastus"
	LocationSecondary = "westeurope"
	LocationTernary   = "eastasia"
)

func ModeFromEnv() Mode {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvVarMode))) {
	case "record":
		return ModeRecord
	case "replay":
		return ModeReplay
	default:
		return ModeOff
	}
}

func Enabled() bool {
	return ModeFromEnv() != ModeOff
}

type Recorder struct {
	t    *testing.T
	mode Mode
	rec  *recorder.Recorder
	san  *sanitizer
}

// The active recorder is a global with a read/write lock so that BuildTestClient can access it without needing to pass a *Recorder around
// It's safe to use despite cross-package concurrency because go spawns a new OS process for each test package
var (
	activeRecorder     *Recorder
	activeRecorderLock sync.RWMutex
)

func Active() *Recorder {
	activeRecorderLock.RLock()
	defer activeRecorderLock.RUnlock()
	return activeRecorder
}

func setActiveRecorder(r *Recorder) {
	activeRecorderLock.Lock()
	defer activeRecorderLock.Unlock()
	if activeRecorder != nil {
		r.t.Fatalf("vcr: cannot start recorder for test %q because the recorder for test %q is still active; "+
			"VCR tests must run serially — do not use t.Parallel() or resource.ParallelTest with a VCR test",
			r.t.Name(), activeRecorder.t.Name())
	}
	activeRecorder = r
}

func clearActiveRecorder(r *Recorder) {
	activeRecorderLock.Lock()
	defer activeRecorderLock.Unlock()
	if activeRecorder == r {
		activeRecorder = nil
		return
	}
	other := "(nil)"
	if activeRecorder != nil {
		other = activeRecorder.t.Name()
	}
	r.t.Fatalf("vcr: the active recorder changed out from under test %q before it stopped (active is now %q); "+
		"VCR tests must run serially — do not use t.Parallel() or resource.ParallelTest with a VCR test",
		r.t.Name(), other)
}

func New(t *testing.T) *Recorder {
	t.Helper()

	r := &Recorder{t: t, mode: ModeFromEnv()}
	if r.mode == ModeOff {
		return r
	}

	realSubscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	realTenantID := os.Getenv("ARM_TENANT_ID")

	// The sanitizer redacts known sensitive values and then audits every
	// interaction, failing the recording (fail-closed) if a secret survives.
	san := newSanitizer(realSubscriptionID, realTenantID)
	// Canonicalize non-secret identifiers and redact the credential material the
	// harness itself supplies, so none of it can be written to a cassette.
	san.AddReplacement(os.Getenv("ARM_CLIENT_ID"), CanonicalClientID)
	san.AddReplacement(os.Getenv("ARM_READER_CLIENT_ID"), CanonicalClientID)
	san.AddSecret(os.Getenv("ARM_CLIENT_SECRET"))
	san.AddSecret(os.Getenv("ARM_READER_CLIENT_SECRET"))
	san.AddSecret(os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"))
	san.AddSecret(os.Getenv("ARM_OIDC_TOKEN"))
	r.san = san

	recorderMode := recorder.ModeReplayOnly
	if r.mode == ModeRecord {
		recorderMode = recorder.ModeRecordOnly
	}

	cassetteName := filepath.Join("testdata", "cassettes", strings.ReplaceAll(t.Name(), "/", "_"))
	rec, err := recorder.New(
		cassetteName,
		recorder.WithMode(recorderMode),
		recorder.WithSkipRequestLatency(true),
		recorder.WithMatcher(methodURLMatcher),
		recorder.WithHook(san.BeforeSave, recorder.BeforeSaveHook),
	)
	if err != nil {
		t.Fatalf("creating VCR recorder: %v", err)
	}
	r.rec = rec

	if r.mode == ModeReplay {
		t.Setenv("TF_ACC", "1")
		t.Setenv("ARM_CLIENT_ID", CanonicalClientID)
		t.Setenv("ARM_CLIENT_SECRET", CanonicalClientSecret)
		t.Setenv("ARM_TENANT_ID", CanonicalTenantID)
		t.Setenv("ARM_SUBSCRIPTION_ID", CanonicalSubscriptionID)
		t.Setenv("ARM_USE_CLI", "false")
		t.Setenv("ARM_USE_MSI", "false")
		t.Setenv("ARM_USE_OIDC", "false")
		t.Setenv("ARM_USE_AKS_WORKLOAD_IDENTITY", "false")
	}

	setActiveRecorder(r)
	return r
}

func (r *Recorder) Enabled() bool {
	return r.mode != ModeOff
}

func (r *Recorder) Mode() Mode {
	return r.mode
}

// AddSecret registers a sensitive value discovered while a test runs (for
// example a storage account key or password returned by Azure) so that it is
// scrubbed from the cassette before it is saved. It is a no-op when the
// recorder is not recording, so tests may call it unconditionally.
func (r *Recorder) AddSecret(v string) {
	if r == nil || r.san == nil {
		return
	}
	r.san.AddSecret(v)
}

// RegisterSecret registers a sensitive value with the active recorder so it is
// scrubbed from the cassette before it is saved. It is safe to call
// unconditionally: when VCR is disabled (or no recorder is active) it does
// nothing. Use it from tests that read secret values back from Azure, e.g.
// vcr.RegisterSecret(storageAccountKey).
func RegisterSecret(v string) {
	Active().AddSecret(v)
}

func (r *Recorder) Stop() {
	if r.rec == nil {
		return
	}
	clearActiveRecorder(r)
	if err := r.rec.Stop(); err != nil {
		r.t.Errorf("stopping VCR recorder: %v", err)
	}
}

func (r *Recorder) Transport() policy.Transporter {
	if r.rec == nil {
		return nil
	}
	return r.rec.GetDefaultClient()
}

func (r *Recorder) Credential() azcore.TokenCredential {
	if r.mode == ModeReplay {
		return fakeCredential{}
	}
	return nil
}

func (r *Recorder) SubscriptionID() string {
	if r.mode == ModeReplay {
		return CanonicalSubscriptionID
	}
	return os.Getenv("ARM_SUBSCRIPTION_ID")
}

func (r *Recorder) TenantID() string {
	if r.mode == ModeReplay {
		return CanonicalTenantID
	}
	return os.Getenv("ARM_TENANT_ID")
}

func (r *Recorder) ConfigureOption(o *clients.Option) {
	if r.rec == nil {
		return
	}
	o.Transport = r.Transport()
	if r.mode == ModeReplay {
		o.Cred = fakeCredential{}
		o.SubscriptionId = CanonicalSubscriptionID
		o.TenantId = CanonicalTenantID
	}
}

func methodURLMatcher(r *http.Request, i cassette.Request) bool {
	return r.Method == i.Method && r.URL.String() == i.URL
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

func RandInt(name string) int {
	r := deterministicRand(name + ":int")
	return int(r.Int63n(900000000000000000) + 100000000000000000)
}

func RandString(name string, length int) string {
	r := deterministicRand(name + ":str")
	b := make([]byte, length)
	for i := range b {
		b[i] = charSetAlphaNum[r.Intn(len(charSetAlphaNum))]
	}
	return string(b)
}
