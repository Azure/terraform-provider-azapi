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
			resourceID:  "https://example.com/versions/1",
			action:      "startPendingUpload",
			apiVersion:  "2025-05-01",
			wantPath:    "/versions/1/startPendingUpload",
			wantVersion: "2025-05-01",
		},
		{
			name:        "resource action without suffix",
			resourceID:  "https://example.com/versions/1",
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

			if strings.Contains(rawURL, "/api-version=") {
				t.Fatalf("API version was appended as a path segment: %s", rawURL)
			}
		})
	}
}
