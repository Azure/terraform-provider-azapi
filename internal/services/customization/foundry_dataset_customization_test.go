package customization

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func TestDatasetUploadDetails(t *testing.T) {
	response := map[string]interface{}{
		"blobReference": map[string]interface{}{
			"blobUri": "https://storage.blob.core.windows.net/container",
			"credential": map[string]interface{}{
				"type":   "SAS",
				"sasUri": "https://storage.blob.core.windows.net/container?sr=c&sig=test",
			},
		},
		"blobReferenceForConsumption": map[string]interface{}{
			"blobUri": "https://storage.blob.core.windows.net/container",
		},
	}

	uploadURL, dataURI, err := datasetUploadDetails(response)
	if err != nil {
		t.Fatalf("datasetUploadDetails returned an error: %v", err)
	}
	if uploadURL != "https://storage.blob.core.windows.net/container?sr=c&sig=test" {
		t.Fatalf("unexpected upload URL: %q", uploadURL)
	}
	if dataURI != "https://storage.blob.core.windows.net/container" {
		t.Fatalf("unexpected data URI: %q", dataURI)
	}
}

func TestDatasetSourceInfo(t *testing.T) {
	t.Run("without checksum", func(t *testing.T) {
		sourceURL, checksum, verify, err := datasetSourceInfo(map[string]interface{}{
			"source_url": "https://example.com/data.jsonl",
		})
		if err != nil {
			t.Fatalf("datasetSourceInfo returned an error: %v", err)
		}
		if sourceURL != "https://example.com/data.jsonl" || checksum != "" || verify {
			t.Fatalf("unexpected result: %q, %q, %t", sourceURL, checksum, verify)
		}
	})

	t.Run("normalizes checksum", func(t *testing.T) {
		sourceURL, checksum, verify, err := datasetSourceInfo(map[string]interface{}{
			"source_url":    "https://example.com/data.jsonl",
			"source_sha256": strings.Repeat("AB", sha256.Size),
		})
		if err != nil {
			t.Fatalf("datasetSourceInfo returned an error: %v", err)
		}
		if sourceURL != "https://example.com/data.jsonl" ||
			checksum != strings.Repeat("ab", sha256.Size) ||
			!verify {
			t.Fatalf("unexpected result: %q, %q, %t", sourceURL, checksum, verify)
		}
	})

	t.Run("rejects invalid checksum", func(t *testing.T) {
		_, _, _, err := datasetSourceInfo(map[string]interface{}{
			"source_url":    "https://example.com/data.jsonl",
			"source_sha256": "not-a-sha256",
		})
		if err == nil {
			t.Fatal("expected invalid source_sha256 to return an error")
		}
	})

	t.Run("requires source URL", func(t *testing.T) {
		_, _, _, err := datasetSourceInfo(map[string]interface{}{})
		if err == nil {
			t.Fatal("expected missing source_url to return an error")
		}
	})
}

func TestDatasetVersionRequestBody(t *testing.T) {
	body, version, err := datasetVersionRequestBody(map[string]interface{}{
		"name":        "example-dataset",
		"description": "example",
		"format":      "jsonl",
		"source_url":  "https://example.com/data.jsonl",
	}, "1", "https://storage.blob.core.windows.net/container")
	if err != nil {
		t.Fatalf("datasetVersionRequestBody returned an error: %v", err)
	}
	if version != "1" {
		t.Fatalf("unexpected version: %q", version)
	}

	expected := map[string]interface{}{
		"name":        "example-dataset",
		"version":     "1",
		"description": "example",
		"type":        "uri-file",
		"dataUri":     "https://storage.blob.core.windows.net/container",
		"format":      "jsonl",
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected version request body:\n got: %#v\nwant: %#v", body, expected)
	}
}

func TestDatasetPendingUploadBody(t *testing.T) {
	body, err := datasetPendingUploadBody(map[string]interface{}{
		"pending_upload_id": "pending-id",
		"connection_name":   "connection",
	})
	if err != nil {
		t.Fatalf("datasetPendingUploadBody returned an error: %v", err)
	}

	expected := map[string]interface{}{
		"pendingUploadType": "BlobReference",
		"pendingUploadId":   "pending-id",
		"connectionName":    "connection",
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected pending upload body:\n got: %#v\nwant: %#v", body, expected)
	}
}

func TestDatasetBlobUploadURL(t *testing.T) {
	t.Run("container SAS", func(t *testing.T) {
		got, err := datasetBlobUploadURL(
			"https://storage.blob.core.windows.net/container?sr=c&sig=test",
			"data.jsonl",
		)
		if err != nil {
			t.Fatalf("datasetBlobUploadURL returned an error: %v", err)
		}
		want := "https://storage.blob.core.windows.net/container/data.jsonl?sr=c&sig=test"
		if got != want {
			t.Fatalf("unexpected upload URL: %q", got)
		}
	})

	t.Run("blob SAS", func(t *testing.T) {
		want := "https://storage.blob.core.windows.net/container/data.jsonl?sr=b&sig=test"
		got, err := datasetBlobUploadURL(want, "ignored.jsonl")
		if err != nil {
			t.Fatalf("datasetBlobUploadURL returned an error: %v", err)
		}
		if got != want {
			t.Fatalf("unexpected upload URL: %q", got)
		}
	})
}

func TestStreamDatasetToUpload(t *testing.T) {
	const contents = "message\nhello\n"

	sourceServer := httptest.NewServer(http.HandlerFunc(func(
		response http.ResponseWriter,
		request *http.Request,
	) {
		if request.Method != http.MethodGet || request.URL.Path != "/data.jsonl" {
			t.Errorf("unexpected source request: %s %s", request.Method, request.URL)
		}
		_, _ = io.WriteString(response, contents)
	}))
	defer sourceServer.Close()

	uploadServer := httptest.NewServer(http.HandlerFunc(func(
		response http.ResponseWriter,
		request *http.Request,
	) {
		if request.Method != http.MethodPut ||
			request.URL.Path != "/container/data.jsonl" {
			t.Errorf("unexpected upload request: %s %s", request.Method, request.URL)
		}
		if request.Header.Get("Content-Type") != "application/octet-stream" {
			t.Errorf("unexpected content type: %q", request.Header.Get("Content-Type"))
		}
		if request.Header.Get("x-ms-blob-type") != "BlockBlob" {
			t.Errorf("unexpected blob type: %q", request.Header.Get("x-ms-blob-type"))
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Errorf("reading upload body: %v", err)
		}
		if string(body) != contents {
			t.Errorf("unexpected upload body: %q", body)
		}
		response.WriteHeader(http.StatusCreated)
	}))
	defer uploadServer.Close()

	sum := sha256.Sum256([]byte(contents))
	expectedSHA256 := hex.EncodeToString(sum[:])
	computedSHA256, err := streamDatasetToUpload(
		t.Context(),
		sourceServer.URL+"/data.jsonl",
		uploadServer.URL+"/container?sr=c&sig=test",
		expectedSHA256,
		true,
	)
	if err != nil {
		t.Fatalf("streamDatasetToUpload returned an error: %v", err)
	}
	if computedSHA256 != expectedSHA256 {
		t.Fatalf("unexpected computed SHA-256: %q", computedSHA256)
	}
}

func TestDatasetHTTPClientsRefuseRedirects(t *testing.T) {
	redirectServer := httptest.NewServer(http.HandlerFunc(func(
		response http.ResponseWriter,
		_ *http.Request,
	) {
		response.Header().Set("Location", "https://example.com/final")
		response.WriteHeader(http.StatusFound)
	}))
	defer redirectServer.Close()

	for name, client := range map[string]*http.Client{
		"source": datasetSourceHTTPClient(),
		"upload": datasetUploadHTTPClient(),
	} {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodGet,
				redirectServer.URL,
				nil,
			)
			if err != nil {
				t.Fatalf("creating request: %v", err)
			}

			response, err := client.Do(request)
			if response != nil {
				response.Body.Close()
			}
			if err == nil {
				t.Fatal("expected redirect to be refused")
			}
		})
	}
}

func TestFoundryDatasetReadOutput(t *testing.T) {
	t.Run("preserves response checksum", func(t *testing.T) {
		customization := FoundryDatasetCustomization{}
		output, err := customization.AugmentReadOutput(
			map[string]interface{}{
				"computed_sha256": "abc",
			},
			map[string]interface{}{},
		)
		if err != nil {
			t.Fatalf("AugmentReadOutput returned an error: %v", err)
		}
		values, ok := output.(map[string]interface{})
		if !ok || values["computed_sha256"] != "abc" {
			t.Fatalf("computed_sha256 was not preserved in output: %#v", output)
		}
	})

	t.Run("computes checksum from source URL", func(t *testing.T) {
		const contents = "dataset"
		sourceServer := httptest.NewServer(http.HandlerFunc(func(
			response http.ResponseWriter,
			_ *http.Request,
		) {
			_, _ = io.WriteString(response, contents)
		}))
		defer sourceServer.Close()

		sum := sha256.Sum256([]byte(contents))
		expectedSHA256 := hex.EncodeToString(sum[:])
		customization := FoundryDatasetCustomization{}
		output, err := customization.AugmentReadOutput(
			map[string]interface{}{"name": "example-dataset"},
			map[string]interface{}{"source_url": sourceServer.URL + "/data.jsonl"},
		)
		if err != nil {
			t.Fatalf("AugmentReadOutput returned an error: %v", err)
		}
		values, ok := output.(map[string]interface{})
		if !ok || values["computed_sha256"] != expectedSHA256 {
			t.Fatalf("unexpected output: %#v", output)
		}
	})
}

func TestFoundryDatasetPlanBody(t *testing.T) {
	customization := FoundryDatasetCustomization{}
	planBody, err := customization.PlanBodyFunc()(
		map[string]interface{}{
			"name":   "example-dataset",
			"format": "jsonl",
		},
		map[string]interface{}{
			"version":         "1",
			"type":            "uri-file",
			"format":          "jsonl",
			"computed_sha256": "abc",
		},
	)
	if err != nil {
		t.Fatalf("PlanBodyFunc returned an error: %v", err)
	}

	values, ok := planBody.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected plan body type: %#v", planBody)
	}
	if values["version"] != "1" || values["type"] != "uri-file" {
		t.Fatalf("state defaults were not copied into plan body: %#v", values)
	}
	if _, exists := values["computed_sha256"]; exists {
		t.Fatalf("computed_sha256 must not be copied into plan body: %#v", values)
	}
}

func TestDatasetDefaults(t *testing.T) {
	t.Run("sets defaults", func(t *testing.T) {
		body := map[string]interface{}{
			"format": "jsonl",
		}

		if err := setDatasetDefaults(body, "1"); err != nil {
			t.Fatalf("setDatasetDefaults returned an error: %v", err)
		}

		if body["version"] != "1" {
			t.Fatalf("unexpected version: %#v", body["version"])
		}
		if body["type"] != "uri-file" {
			t.Fatalf("unexpected type: %#v", body["type"])
		}
	})

	t.Run("requires format", func(t *testing.T) {
		err := setDatasetDefaults(
			map[string]interface{}{},
			"1",
		)
		if err == nil {
			t.Fatal("expected missing format to return an error")
		}
	})

	t.Run("rejects generated version", func(t *testing.T) {
		err := setDatasetDefaults(
			map[string]interface{}{},
			"__generated__",
		)
		if err == nil {
			t.Fatal("expected generated version to return an error")
		}
	})

	t.Run("allows folder uploads", func(t *testing.T) {
		body := map[string]interface{}{
			"format": "jsonl",
			"type":   "uri-folder",
		}

		if err := setDatasetDefaults(
			body,
			"1",
		); err != nil {
			t.Fatalf("setDatasetDefaults returned an error: %v", err)
		}
		if body["type"] != "uri-folder" {
			t.Fatalf("unexpected type: %#v", body["type"])
		}
	})
}

func TestDatasetVersionRequestBodyRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		body map[string]interface{}
	}{
		{
			name: "missing name",
			body: map[string]interface{}{
				"description": "example",
				"format":      "jsonl",
			},
		},
		{
			name: "missing description",
			body: map[string]interface{}{
				"name":   "example",
				"format": "jsonl",
			},
		},
		{
			name: "missing format",
			body: map[string]interface{}{
				"name":        "example",
				"description": "example",
			},
		},
		{
			name: "invalid type",
			body: map[string]interface{}{
				"name":        "example",
				"description": "example",
				"format":      "jsonl",
				"type":        "invalid",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := datasetVersionRequestBody(
				test.body,
				"1",
				"https://storage.example/data.jsonl",
			)
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestDatasetVersionRequestBodyAllowsFolderType(t *testing.T) {
	body, _, err := datasetVersionRequestBody(map[string]interface{}{
		"name":        "example-dataset",
		"description": "example",
		"format":      "jsonl",
		"type":        "uri-folder",
	}, "1", "https://storage.blob.core.windows.net/container")
	if err != nil {
		t.Fatalf("datasetVersionRequestBody returned an error: %v", err)
	}
	if body["type"] != "uri-folder" {
		t.Fatalf("unexpected dataset type: %#v", body["type"])
	}
}

func TestFoundryDatasetCustomizationLifecycle(t *testing.T) {
	customization := FoundryDatasetCustomization{}

	if customization.GetResourceType() != "Microsoft.Foundry/datasets/versions" {
		t.Fatalf("unexpected resource type: %q", customization.GetResourceType())
	}
	if customization.CreateFunc() == nil ||
		customization.ReadFunc() == nil ||
		customization.UpdateFunc() == nil {
		t.Fatal("dataset customization must define create, read, and update behavior")
	}
	if customization.DeleteFunc() != nil {
		t.Fatal("dataset customization should use the generic delete behavior")
	}

	err := customization.UpdateFunc()(
		t.Context(),
		clients.Client{},
		parse.DataPlaneResourceId{},
		nil,
		clients.RequestOptions{},
	)
	if err == nil || !strings.Contains(err.Error(), "immutable") {
		t.Fatalf("unexpected immutable update error: %v", err)
	}
}

func TestDatasetUploadDetailsRejectsMalformedResponses(t *testing.T) {
	tests := []struct {
		name     string
		response interface{}
	}{
		{
			name:     "nil response",
			response: nil,
		},
		{
			name:     "missing blob reference",
			response: map[string]interface{}{},
		},
		{
			name: "missing credential",
			response: map[string]interface{}{
				"blobReference": map[string]interface{}{},
			},
		},
		{
			name: "missing SAS URI",
			response: map[string]interface{}{
				"blobReference": map[string]interface{}{
					"credential": map[string]interface{}{},
				},
			},
		},
		{
			name: "missing blob URI",
			response: map[string]interface{}{
				"blobReference": map[string]interface{}{
					"credential": map[string]interface{}{
						"sasUri": "https://storage.example/container?sr=c",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := datasetUploadDetails(test.response)
			if err == nil {
				t.Fatal("expected malformed response to return an error")
			}
		})
	}
}

func TestStreamDatasetToUploadRejectsChecksumMismatch(
	t *testing.T,
) {
	const contents = "dataset contents"

	sourceServer := httptest.NewServer(http.HandlerFunc(
		func(response http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(response, contents)
		},
	))
	defer sourceServer.Close()

	uploadCalled := false
	uploadServer := httptest.NewServer(http.HandlerFunc(
		func(response http.ResponseWriter, _ *http.Request) {
			uploadCalled = true
			response.WriteHeader(http.StatusCreated)
		},
	))
	defer uploadServer.Close()

	_, err := streamDatasetToUpload(
		t.Context(),
		sourceServer.URL+"/dataset.jsonl",
		uploadServer.URL+"/container?sr=c&sig=test",
		strings.Repeat("0", sha256.Size*2),
		true,
	)
	if err == nil {
		t.Fatal("expected checksum mismatch error")
	}

	if !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Fatalf("unexpected error: %v", err)
	}

	if uploadCalled {
		t.Fatal("upload must not occur after checksum mismatch")
	}
}

func TestStreamDatasetToUploadRejectsUploadFailure(t *testing.T) {
	const contents = "dataset contents"

	sourceServer := httptest.NewServer(http.HandlerFunc(
		func(response http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(response, contents)
		},
	))
	defer sourceServer.Close()

	uploadServer := httptest.NewServer(http.HandlerFunc(
		func(response http.ResponseWriter, _ *http.Request) {
			response.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer uploadServer.Close()

	sum := sha256.Sum256([]byte(contents))
	checksum := hex.EncodeToString(sum[:])

	_, err := streamDatasetToUpload(
		t.Context(),
		sourceServer.URL+"/dataset.jsonl",
		uploadServer.URL+"/container?sr=c&sig=test",
		checksum,
		true,
	)
	if err == nil {
		t.Fatal("expected upload failure")
	}

	if !strings.Contains(err.Error(), "uploading dataset returned HTTP") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDatasetSafeErrorRedactsCause(t *testing.T) {
	const sensitiveValue = "https://storage.example/container?sig=secret"

	actualError := errors.New(
		"request failed for " + sensitiveValue,
	)

	redactedError := datasetSafeError(
		"uploading dataset",
		actualError,
	)

	const expectedMessage = "uploading dataset: request failed"

	if got := redactedError.Error(); got != expectedMessage {
		t.Fatalf("unexpected redacted error: got %q, want %q", got, expectedMessage)
	}

	if strings.Contains(redactedError.Error(), sensitiveValue) {
		t.Fatal("redacted error contains sensitive information")
	}

	if errors.Is(redactedError, actualError) {
		t.Fatal("original error must not be recoverable")
	}

	if errors.Unwrap(redactedError) != nil {
		t.Fatal("redacted error must not unwrap to the original error")
	}

	var typedError datasetRedactedError
	if !errors.As(redactedError, &typedError) {
		t.Fatal("expected datasetRedactedError")
	}

	if typedError.operation != "uploading dataset" {
		t.Fatalf("unexpected operation: %q", typedError.operation)
	}
}
