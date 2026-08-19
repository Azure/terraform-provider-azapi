package clients

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildDataPlaneActionURL(t *testing.T) {
	tests := []struct {
		name        string
		resourceID  string
		action      string
		apiVersion  string
		wantPath    string
		wantVersion string
	}{
		{
			name:        "action uses query parameter",
			resourceID:  "example.com/versions/1",
			action:      "startPendingUpload",
			apiVersion:  "2025-05-01",
			wantPath:    "/versions/1/startPendingUpload",
			wantVersion: "2025-05-01",
		},
		{
			name:        "resource action preserves query parameters",
			resourceID:  "example.com/versions/1?foo=bar",
			action:      "",
			apiVersion:  "2025-05-01",
			wantPath:    "/versions/1",
			wantVersion: "2025-05-01",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rawURL, err := buildDataPlaneActionURL(
				test.resourceID,
				test.action,
				test.apiVersion,
			)
			if err != nil {
				t.Fatalf("buildDataPlaneActionURL returned an error: %v", err)
			}

			parsedURL, err := url.Parse(rawURL)
			if err != nil {
				t.Fatalf("parsing generated URL: %v", err)
			}

			if parsedURL.Path != test.wantPath {
				t.Fatalf(
					"unexpected path: got %q, want %q",
					parsedURL.Path,
					test.wantPath,
				)
			}

			if parsedURL.Query().Get("api-version") != test.wantVersion {
				t.Fatalf(
					"unexpected API version: got %q, want %q",
					parsedURL.Query().Get("api-version"),
					test.wantVersion,
				)
			}

			if parsedURL.Scheme != "https" {
				t.Fatalf("unexpected scheme: got %q, want %q", parsedURL.Scheme, "https")
			}
			if parsedURL.Host != "example.com" {
				t.Fatalf("unexpected host: got %q, want %q", parsedURL.Host, "example.com")
			}
			if test.name == "resource action preserves query parameters" &&
				parsedURL.Query().Get("foo") != "bar" {
				t.Fatalf("existing query parameter was not preserved: %s", rawURL)
			}
			if strings.Contains(parsedURL.Path, "api-version") {
				t.Fatalf("API version was appended as a path segment: %s", rawURL)
			}
		})
	}
}
