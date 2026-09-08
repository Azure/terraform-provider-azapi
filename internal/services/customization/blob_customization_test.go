package customization

import (
	"encoding/base64"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestBuildBlobRequest(t *testing.T) {
	content := []byte{0, 1, 2, 255}
	body := map[string]interface{}{
		"content_base64":      base64.StdEncoding.EncodeToString(content),
		"content_type":        "application/octet-stream",
		"content_disposition": "attachment; filename=blob.bin",
		"metadata": map[string]interface{}{
			"environment": "test",
		},
		"access_tier": "Cool",
	}
	actualContent, headers, err := buildBlobRequest(body, "2023-11-03", map[string]string{"If-Match": `"etag"`})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actualContent, content) {
		t.Fatalf("unexpected content %v", actualContent)
	}
	expectedHeaders := map[string]string{
		"Content-Type":                  "application/octet-stream",
		"x-ms-blob-content-disposition": "attachment; filename=blob.bin",
		"x-ms-meta-environment":         "test",
		"x-ms-access-tier":              "Cool",
		"If-Match":                      `"etag"`,
		"x-ms-blob-type":                "BlockBlob",
	}
	if !reflect.DeepEqual(headers, expectedHeaders) {
		t.Fatalf("unexpected headers %#v", headers)
	}
}

func TestBuildBlobRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		version string
		body    map[string]interface{}
		headers map[string]string
		want    string
	}{
		{
			name:    "invalid base64",
			version: "2023-11-03",
			body:    map[string]interface{}{"content_base64": "***"},
			want:    "invalid content_base64",
		},
		{
			name:    "unknown field",
			version: "2023-11-03",
			body:    map[string]interface{}{"content_base64": "", "contents": "value"},
			want:    `unsupported blob body field "contents"`,
		},
		{
			name:    "old version",
			version: "2017-07-29",
			body:    map[string]interface{}{"content_base64": ""},
			want:    "does not support Microsoft Entra authentication",
		},
		{
			name:    "reserved header",
			version: "2023-11-03",
			body:    map[string]interface{}{"content_base64": ""},
			headers: map[string]string{"Content-Length": "3"},
			want:    "cannot be overridden",
		},
		{
			name:    "invalid tier",
			version: "2023-11-03",
			body:    map[string]interface{}{"content_base64": "", "access_tier": "Premium"},
			want:    "must be one of",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := buildBlobRequest(test.body, test.version, test.headers)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestBlobProperties(t *testing.T) {
	headers := http.Header{
		"Content-Type":          []string{"text/plain"},
		"Content-Length":        []string{"31"},
		"Content-Md5":           []string{"YWJj"},
		"Etag":                  []string{`"0x8"`},
		"Last-Modified":         []string{"Mon, 08 Sep 2026 01:02:03 GMT"},
		"X-Ms-Access-Tier":      []string{"Hot"},
		"X-Ms-Blob-Type":        []string{"BlockBlob"},
		"X-Ms-Meta-Environment": []string{"test"},
	}
	actual, err := blobProperties(headers)
	if err != nil {
		t.Fatal(err)
	}
	if actual["content_length"] != int64(31) || actual["content_type"] != "text/plain" || actual["content_encoding"] != nil {
		t.Fatalf("unexpected properties %#v", actual)
	}
	metadata, ok := actual["metadata"].(map[string]interface{})
	if !ok || metadata["environment"] != "test" {
		t.Fatalf("unexpected metadata %#v", actual["metadata"])
	}
}
