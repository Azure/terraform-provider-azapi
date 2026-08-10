package vcr

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

func sanitizeInteraction(realSubscriptionID, realTenantID string) recorder.HookFunc {
	replacements := map[string]string{}
	if realSubscriptionID != "" {
		replacements[realSubscriptionID] = CanonicalSubscriptionID
	}
	if realTenantID != "" {
		replacements[realTenantID] = CanonicalTenantID
	}

	redact := func(s string) string {
		for from, to := range replacements {
			s = strings.ReplaceAll(s, from, to)
		}
		return s
	}
	redactQueryParameters := func(s string) string {
		return queryParameterPattern.ReplaceAllString(s, `${1}`+redactedValue)
	}
	redactResponseBody := func(s string) string {
		s = redact(s)
		s = etagJSONPattern.ReplaceAllString(s, `${1}"W/\"`+CanonicalSubscriptionID+`\""`)
		return resourceGUIDJSONPattern.ReplaceAllString(s, `${1}"`+CanonicalSubscriptionID+`"`)
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

	return func(i *cassette.Interaction) error {
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
		return nil
	}
}
