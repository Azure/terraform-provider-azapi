package vcr

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

// The audit is a fail-closed backstop for sanitizeInteraction's denylist. After
// redaction it re-scans every interaction for high-signal secret shapes and, if
// any survive, returns an error so go-vcr refuses to persist the cassette. The
// detectors are intentionally high-precision (few false positives): the goal is
// to catch sensitive material the denylist does not yet know about, not to
// replace it. When a legitimate recording trips the audit, redact the value in
// sanitize_interaction.go or register it at runtime with vcr.RegisterSecret.

// secretDetector is a single high-precision rule used by auditInteraction to
// recognize a class of value that must never be written to a cassette.
type secretDetector struct {
	name string
	re   *regexp.Regexp
	// valueGroup is the submatch index holding the sensitive value; 0 means the
	// whole match. A match whose value is a known-safe placeholder is ignored.
	valueGroup int
	// labelGroup, when > 0, is a submatch index (e.g. a JSON field name) added
	// to the finding for context. It is never the secret value itself.
	labelGroup int
	// predicate, when set, must return true for a match to count as a finding.
	predicate func(value string) bool
}

var (
	// A JWT/bearer token: three dot-separated base64url segments beginning with
	// the "eyJ" that base64-encodes `{"`. Opaque continuation tokens lack the
	// dotted structure, so this stays specific to real tokens.
	jwtPattern = regexp.MustCompile(`eyJ[A-Za-z0-9_-]{6,}\.[A-Za-z0-9_-]{6,}\.[A-Za-z0-9_-]*`)
	// A PEM private key block of any algorithm (RSA, EC, OPENSSH, ...).
	pemPrivateKeyPattern = regexp.MustCompile(`-----BEGIN (?:[A-Z0-9]+ )*PRIVATE KEY-----`)
	// A shared-access-signature `sig=` parameter, as found in Storage/Service Bus SAS URLs.
	sasSignaturePattern = regexp.MustCompile(`(?i)[?&;]sig=([A-Za-z0-9%+/=_-]{10,})`)
	// A `SharedAccessSignature=...` assignment (Service Bus / Event Hub / IoT Hub).
	sasTokenPattern = regexp.MustCompile(`(?i)SharedAccessSignature\s*[=:]\s*(\S{6,})`)
	// The secret assignment inside a connection string.
	connStringKeyPattern = regexp.MustCompile(`(?i)\b(?:AccountKey|SharedAccessKey|AccessKey|Password|Pwd)=([^;"'\s]{6,})`)
	// JSON fields whose name is, in Azure APIs, essentially always a secret.
	secretFieldPattern = regexp.MustCompile(`(?i)"(primaryKey|secondaryKey|primaryMasterKey|secondaryMasterKey|primaryReadonlyMasterKey|secondaryReadonlyMasterKey|sasToken|accountSasToken|serviceSasToken|accountKey|sharedAccessKey|accessKey|storageAccountKey|primaryAccessKey|secondaryAccessKey|aadClientSecret|clientSecret|servicePrincipalClientSecret|adminPassword|administratorPassword|administratorLoginPassword|password|passwordValue|secret|secretValue|connectionString|primaryConnectionString|secondaryConnectionString)"\s*:\s*"((?:\\.|[^"\\])*)"`)
	// A `"value":"<base64/hex blob>"` field (e.g. a storage account list-keys
	// response). Gated by an entropy check so ordinary values do not trip it.
	valueFieldPattern = regexp.MustCompile(`(?i)"(?:value|keyData|secretText|contentBytes)"\s*:\s*"([A-Za-z0-9+/]{32,}={0,2})"`)

	secretDetectors = []secretDetector{
		{name: "JWT/bearer token", re: jwtPattern},
		{name: "PEM private key", re: pemPrivateKeyPattern},
		{name: "SAS signature", re: sasSignaturePattern, valueGroup: 1},
		{name: "shared access signature", re: sasTokenPattern, valueGroup: 1},
		{name: "connection-string secret", re: connStringKeyPattern, valueGroup: 1},
		{name: "secret JSON field", re: secretFieldPattern, valueGroup: 2, labelGroup: 1},
		{name: "high-entropy value field", re: valueFieldPattern, valueGroup: 1, predicate: isHighEntropy},
	}
)

// auditInteraction re-scans an already-sanitized interaction for residual
// secrets. It returns a non-nil error listing every finding (never the secret
// value itself) when something slips past sanitizeInteraction's denylist. As
// the second half of (*sanitizer).BeforeSave, a non-nil result makes go-vcr
// refuse to persist the cassette.
func auditInteraction(i *cassette.Interaction, replacements map[string]string) error {
	safe := safeValues(replacements)

	type surface struct {
		where string
		text  string
	}
	var surfaces []surface
	add := func(where, text string) {
		if text != "" {
			surfaces = append(surfaces, surface{where: where, text: text})
		}
	}

	add("request URL", i.Request.URL)
	add("request URI", i.Request.RequestURI)
	add("request host", i.Request.Host)
	add("request body", i.Request.Body)
	for key, values := range i.Request.Headers {
		for _, v := range values {
			add("request header "+key, v)
		}
	}
	for key, values := range i.Request.Form {
		for _, v := range values {
			add("request form "+key, v)
		}
	}
	add("response body", i.Response.Body)
	for key, values := range i.Response.Headers {
		for _, v := range values {
			add("response header "+key, v)
		}
	}

	findings := map[string]struct{}{}
	for _, s := range surfaces {
		for _, d := range secretDetectors {
			for _, m := range d.re.FindAllStringSubmatch(s.text, -1) {
				value := m[0]
				if d.valueGroup > 0 && d.valueGroup < len(m) {
					value = m[d.valueGroup]
				}
				if isSafeValue(value, safe) {
					continue
				}
				if d.predicate != nil && !d.predicate(value) {
					continue
				}
				desc := d.name
				if d.labelGroup > 0 && d.labelGroup < len(m) {
					desc = fmt.Sprintf("%s %q", d.name, m[d.labelGroup])
				}
				findings[fmt.Sprintf("%s: %s", s.where, desc)] = struct{}{}
			}
		}

		// Residue check: any known-sensitive value that should have been
		// rewritten but is still present (e.g. because Azure echoed it back in a
		// different letter case that ReplaceAll missed).
		lowerText := strings.ToLower(s.text)
		for from, to := range replacements {
			if from == "" || from == to {
				continue
			}
			if strings.Contains(lowerText, strings.ToLower(from)) {
				findings[fmt.Sprintf("%s: unredacted value (expected %q)", s.where, to)] = struct{}{}
			}
		}
	}

	if len(findings) == 0 {
		return nil
	}

	sorted := make([]string, 0, len(findings))
	for f := range findings {
		sorted = append(sorted, f)
	}
	sort.Strings(sorted)

	return fmt.Errorf(
		"vcr: refusing to save cassette interaction %d because potential secret(s) survived sanitization:\n  - %s\n"+
			"Redact these in internal/acceptance/vcr/sanitize_interaction.go, or register the value at runtime with vcr.RegisterSecret(...).",
		i.ID, strings.Join(sorted, "\n  - "),
	)
}

// safeValues returns the set of values that, if matched by a detector, are
// known to be non-sensitive placeholders rather than real secrets.
func safeValues(replacements map[string]string) map[string]struct{} {
	safe := map[string]struct{}{
		redactedValue:           {},
		CanonicalSubscriptionID: {},
		CanonicalTenantID:       {},
		CanonicalObjectID:       {},
		CanonicalClientID:       {},
		CanonicalClientSecret:   {},
	}
	for _, to := range replacements {
		safe[to] = struct{}{}
	}
	return safe
}

// isSafeValue reports whether a captured value is a placeholder rather than a
// real secret. It treats registered replacement targets, empty strings and
// common masks (asterisks, "redacted", "null", ...) as safe so that already
// sanitized content never trips the audit.
func isSafeValue(v string, safe map[string]struct{}) bool {
	if _, ok := safe[v]; ok {
		return true
	}
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return true
	}
	if _, ok := safe[trimmed]; ok {
		return true
	}
	// Runs of a single masking character such as "****" or "xxxx".
	for _, mask := range []string{"*", "x", "X", "#", "."} {
		if strings.Trim(trimmed, mask) == "" {
			return true
		}
	}
	switch strings.ToLower(strings.Trim(trimmed, "[]<>()")) {
	case redactedValue, "hidden", "masked", "sanitized", "sanitised", "null", "none", "placeholder":
		return true
	}
	return false
}

// isHighEntropy reports whether a value looks like random key material rather
// than a human-readable identifier. Storage keys and similar base64 blobs score
// well above the threshold; short or structured strings score below it.
func isHighEntropy(v string) bool {
	if len(v) < 32 {
		return false
	}
	return shannonEntropy(v) >= 4.0
}

// shannonEntropy returns the Shannon entropy, in bits per byte, of s.
func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}
	var counts [256]float64
	for i := 0; i < len(s); i++ {
		counts[s[i]]++
	}
	n := float64(len(s))
	var entropy float64
	for _, c := range counts {
		if c == 0 {
			continue
		}
		p := c / n
		entropy -= p * math.Log2(p)
	}
	return entropy
}
