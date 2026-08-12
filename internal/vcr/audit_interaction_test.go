package vcr

import (
	"net/http"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

// A well-mixed base64 blob (entropy ~5.9) standing in for a storage account key.
const fakeStorageKey = "MFyErXpQ0aZ7kTnLd3Rv9BwHc2sUj6IgKmO1PxV4tYbN8uW5zD7aQeRiShJlC0oP2fXtG9wZ8=="

func responseBodyInteraction(body string) *cassette.Interaction {
	i := &cassette.Interaction{}
	i.Response.Body = body
	return i
}

func TestAuditDetectsUnknownSecrets(t *testing.T) {
	cases := []struct {
		name string
		i    *cassette.Interaction
	}{
		{
			name: "JWT in response body",
			i:    responseBodyInteraction(`{"access_token":"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abcDEFghiJKLmnoPQRstuv"}`),
		},
		{
			name: "PEM private key in response body",
			i:    responseBodyInteraction("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKC\n-----END RSA PRIVATE KEY-----"),
		},
		{
			name: "SAS signature in Location header",
			i: func() *cassette.Interaction {
				i := &cassette.Interaction{}
				i.Response.Headers = http.Header{"Location": {"https://x.blob.core.windows.net/c/b?sv=2021-08-06&sig=aB3dEf7Gh9%2BkLmNoPqRsTuVwXyZ0123456789abcd"}}
				return i
			}(),
		},
		{
			name: "Service Bus SharedAccessKey connection string",
			i:    responseBodyInteraction(`{"primaryConnectionString":"Endpoint=sb://x.servicebus.windows.net/;SharedAccessKeyName=root;SharedAccessKey=aB3dEf7Gh9kLmNoPqRsTuVwXyZ0123456789abcdE=;"}`),
		},
		{
			name: "SQL admin password in connection string",
			i:    responseBodyInteraction("Server=tcp:x.database.windows.net;User ID=admin;Password=P@ssw0rd-Str0ng!;Encrypt=true;"),
		},
		{
			name: "storage account key JSON field",
			i:    responseBodyInteraction(`{"keys":[{"keyName":"key1","value":"` + fakeStorageKey + `","permissions":"FULL"}]}`),
		},
		{
			name: "cosmos primaryMasterKey field",
			i:    responseBodyInteraction(`{"primaryMasterKey":"` + fakeStorageKey + `"}`),
		},
		{
			name: "clientSecret field in request body",
			i: func() *cassette.Interaction {
				i := &cassette.Interaction{}
				i.Request.Body = `{"properties":{"clientSecret":"s3cr3t-value-not-redacted"}}`
				return i
			}(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := auditInteraction(tc.i, nil); err == nil {
				t.Fatalf("expected audit to detect a secret, got nil")
			}
		})
	}
}

func TestAuditAllowsSanitizedContent(t *testing.T) {
	cases := []struct {
		name string
		i    *cassette.Interaction
	}{
		{
			name: "already redacted account key",
			i:    responseBodyInteraction(`{"keys":[{"keyName":"key1","value":"` + redactedValue + `"}]}`),
		},
		{
			name: "already redacted password in connection string",
			i:    responseBodyInteraction("Server=x;User ID=admin;Password=" + redactedValue + ";"),
		},
		{
			name: "masked secret field",
			i:    responseBodyInteraction(`{"clientSecret":"****"}`),
		},
		{
			name: "empty secret field",
			i:    responseBodyInteraction(`{"password":""}`),
		},
		{
			name: "canonical guids and ordinary payload",
			i:    responseBodyInteraction(`{"id":"/subscriptions/` + CanonicalSubscriptionID + `/rg","location":"eastus","value":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`),
		},
		{
			name: "resource guid is not a shared access key name",
			i:    responseBodyInteraction(`{"sharedAccessKeyName":"RootManageSharedAccessKey","keyName":"key1"}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := auditInteraction(tc.i, nil); err != nil {
				t.Fatalf("expected clean audit, got: %v", err)
			}
		})
	}
}

func TestAuditDetectsResidualRegisteredValue(t *testing.T) {
	// The value was registered for redaction but a copy in a different letter
	// case slipped past ReplaceAll; the case-insensitive residue check must
	// still catch it.
	i := responseBodyInteraction(`{"note":"leaked SUPERSECRETVALUE here"}`)
	err := auditInteraction(i, map[string]string{"supersecretvalue": redactedValue})
	if err == nil {
		t.Fatal("expected residue of a registered value to be detected")
	}
	if strings.Contains(err.Error(), "supersecretvalue") || strings.Contains(err.Error(), "SUPERSECRETVALUE") {
		t.Errorf("audit error must not leak the secret value: %v", err)
	}
}

func TestBeforeSaveRedactsRegisteredSecret(t *testing.T) {
	s := newSanitizer("", "")
	s.AddSecret(fakeStorageKey)

	i := responseBodyInteraction(`{"keys":[{"keyName":"key1","value":"` + fakeStorageKey + `"}]}`)
	if err := s.BeforeSave(i); err != nil {
		t.Fatalf("BeforeSave should redact the registered secret and pass, got: %v", err)
	}
	if strings.Contains(i.Response.Body, fakeStorageKey) {
		t.Errorf("registered secret was not redacted from body: %q", i.Response.Body)
	}
	if !strings.Contains(i.Response.Body, redactedValue) {
		t.Errorf("expected redaction placeholder in body: %q", i.Response.Body)
	}
}

func TestBeforeSaveFailsClosedOnUnknownSecret(t *testing.T) {
	s := newSanitizer("", "")
	i := responseBodyInteraction(`{"token":"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abcDEFghiJKLmnoPQRstuv"}`)
	if err := s.BeforeSave(i); err == nil {
		t.Fatal("BeforeSave must fail closed when an unknown secret survives sanitization")
	}
}

func TestShannonEntropy(t *testing.T) {
	if e := shannonEntropy(fakeStorageKey); e < 4.0 {
		t.Errorf("high-entropy key scored too low: %.3f", e)
	}
	if e := shannonEntropy(strings.Repeat("a", 40)); e >= 4.0 {
		t.Errorf("repeated string scored too high: %.3f", e)
	}
	if !isHighEntropy(fakeStorageKey) {
		t.Error("isHighEntropy should be true for a key-like blob")
	}
	if isHighEntropy("eastus") {
		t.Error("isHighEntropy should be false for a short identifier")
	}
}
