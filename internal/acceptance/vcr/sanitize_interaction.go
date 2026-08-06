package vcr

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

const redactedValue = "redacted"

var (
	userAgentPIDPattern     = regexp.MustCompile(`\bpid-[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}\b`)
	queryParameterPattern   = regexp.MustCompile(`([?&](?:c|h|s|t)=)[^&]*`)
	etagJSONPattern         = regexp.MustCompile(`(?i)("etag"[[:space:]]*:[[:space:]]*)"(?:\\.|[^"\\])*"`)
	resourceGUIDJSONPattern = regexp.MustCompile(`(?i)("(?:resourceguid|resourceguide)"[[:space:]]*:[[:space:]]*)"(?:\\.|[^"\\])*"`)
)

// sanitizer scrubs sensitive material from recorded interactions before they
// are persisted to a cassette. Its set of exact-string replacements can be
// extended while a test runs (see AddSecret), which is why the replacements are
// read under a lock every time the BeforeSave hook fires.
type sanitizer struct {
	mu           sync.RWMutex
	replacements map[string]string
}

func newSanitizer(realSubscriptionID, realTenantID string) *sanitizer {
	s := &sanitizer{replacements: make(map[string]string)}
	s.AddReplacement(realSubscriptionID, CanonicalSubscriptionID)
	s.AddReplacement(realTenantID, CanonicalTenantID)
	return s
}

// AddReplacement registers an exact string that is rewritten to `to` in every
// saved interaction. An empty `from` is ignored so callers may pass
// possibly-unset environment variables directly.
func (s *sanitizer) AddReplacement(from, to string) {
	if from == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replacements[from] = to
}

// AddSecret registers a value that must never appear in a saved cassette. The
// value is redacted wherever it occurs and, if any trace survives redaction,
// the audit performed by BeforeSave fails the recording. An empty value is a
// no-op, so callers may pass possibly-unset environment variables directly.
func (s *sanitizer) AddSecret(v string) {
	s.AddReplacement(v, redactedValue)
}

func (s *sanitizer) snapshotReplacements() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.replacements))
	for from, to := range s.replacements {
		out[from] = to
	}
	return out
}

// BeforeSave is the go-vcr BeforeSaveHook. It first redacts every known
// sensitive value and then audits the result: if any residual secret is
// detected it returns an error, which makes go-vcr abort the save so a cassette
// is never written with the secret in it (fail-closed).
func (s *sanitizer) BeforeSave(i *cassette.Interaction) error {
	replacements := s.snapshotReplacements()

	redact := func(str string) string {
		for from, to := range replacements {
			str = strings.ReplaceAll(str, from, to)
		}
		return str
	}
	redactQueryParameters := func(str string) string {
		return queryParameterPattern.ReplaceAllString(str, `${1}`+redactedValue)
	}
	redactResponseBody := func(str string) string {
		str = redact(str)
		str = etagJSONPattern.ReplaceAllString(str, `${1}"W/\"`+CanonicalSubscriptionID+`\""`)
		return resourceGUIDJSONPattern.ReplaceAllString(str, `${1}"`+CanonicalSubscriptionID+`"`)
	}
	redactHeaders := func(h http.Header) {
		for _, sensitive := range []string{"Authorization", "Set-Cookie", "Cookie", "X-Ms-Authorization-Auxiliary"} {
			h.Del(sensitive)
		}
		for key, values := range h {
			normalizedKey := strings.ToLower(key)
			if strings.HasSuffix(normalizedKey, "-request-id") {
				delete(h, key)
				continue
			}
			switch normalizedKey {
			case "client-request-id",
				"correlation-id",
				"request-id",
				"traceparent",
				"tracestate",
				"x-azure-ref",
				"x-ms-operation-identifier",
				"x-msedge-ref",
				"x-vss-e2eid":
				delete(h, key)
				continue
			}
			for idx, value := range values {
				values[idx] = redact(value)
				values[idx] = redactQueryParameters(values[idx])
				if normalizedKey == "user-agent" {
					values[idx] = userAgentPIDPattern.ReplaceAllString(values[idx], "pid-"+redactedValue)
				}
			}
			h[key] = values
		}
	}

	i.Request.URL = redactQueryParameters(redact(i.Request.URL))
	i.Request.RequestURI = redactQueryParameters(redact(i.Request.RequestURI))
	for _, key := range []string{"c", "h", "s", "t"} {
		if _, ok := i.Request.Form[key]; ok {
			i.Request.Form.Set(key, redactedValue)
		}
	}
	i.Request.Body = redact(i.Request.Body)
	i.Response.Body = redactResponseBody(i.Response.Body)
	i.Response.ContentLength = int64(len(i.Response.Body))
	if _, ok := i.Response.Headers["Content-Length"]; ok {
		i.Response.Headers.Set("Content-Length", strconv.Itoa(len(i.Response.Body)))
	}
	redactHeaders(i.Request.Headers)
	redactHeaders(i.Response.Headers)

	return auditInteraction(i, replacements)
}

// sanitizeInteraction returns a BeforeSaveHook backed by a fresh sanitizer. It
// is retained for callers and tests that do not need to register secrets while
// a test runs; the recorder uses (*sanitizer).BeforeSave directly so that
// runtime-registered secrets are honored.
func sanitizeInteraction(realSubscriptionID, realTenantID string) recorder.HookFunc {
	return newSanitizer(realSubscriptionID, realTenantID).BeforeSave
}
