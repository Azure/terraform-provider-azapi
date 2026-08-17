package customization

import (
	"reflect"
	"testing"
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
	sourceURL, checksum, verify, err := datasetSourceInfo(map[string]interface{}{
		"source_url": "https://raw.githubusercontent.com/Azure/terraform-provider-azapi/main/data.jsonl",
	})
	if err != nil {
		t.Fatalf("datasetSourceInfo returned an error: %v", err)
	}
	if sourceURL == "" || checksum != "" || verify {
		t.Fatalf("unexpected optional checksum result: %q, %q, %t", sourceURL, checksum, verify)
	}

	_, _, _, err = datasetSourceInfo(map[string]interface{}{
		"source_url":    "https://raw.githubusercontent.com/Azure/terraform-provider-azapi/main/data.jsonl",
		"source_sha256": "not-a-sha256",
	})
	if err == nil {
		t.Fatal("expected invalid source_sha256 to return an error")
	}

	if err := validateDatasetSourceURL("https://example.com/data.jsonl"); err == nil {
		t.Fatal("expected an unsupported source host to return an error")
	}
}

func TestDatasetVersionRequestBodyDefaults(t *testing.T) {
	body, version, err := datasetVersionRequestBody(map[string]interface{}{
		"name":        "example-dataset",
		"description": "example",
	}, "__generated__", "https://storage.blob.core.windows.net/container")
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
		"type":        "uri_file",
		"dataUri":     "https://storage.blob.core.windows.net/container",
		"format":      "jsonl",
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected version request body:\n got: %#v\nwant: %#v", body, expected)
	}
}

func TestDatasetBlobUploadURL(t *testing.T) {
	got, err := datasetBlobUploadURL("https://storage.blob.core.windows.net/container?sr=c&sig=test", "data.jsonl")
	if err != nil {
		t.Fatalf("datasetBlobUploadURL returned an error: %v", err)
	}
	want := "https://storage.blob.core.windows.net/container/data.jsonl?sr=c&sig=test"
	if got != want {
		t.Fatalf("unexpected blob upload URL: %q", got)
	}
}

func TestFoundryDatasetReadOutput(t *testing.T) {
	customization := FoundryDatasetCustomization{}
	output, err := customization.AugmentReadOutput(
		map[string]interface{}{
			"blobReference": map[string]interface{}{
				"blobUri": "https://storage.blob.core.windows.net/container",
			},
		},
		map[string]interface{}{
			"computed_sha256": "abc",
		},
	)
	if err != nil {
		t.Fatalf("AugmentReadOutput returned an error: %v", err)
	}
	values, ok := output.(map[string]interface{})
	if !ok || values["computed_sha256"] != "abc" {
		t.Fatalf("computed_sha256 was not added to output: %#v", output)
	}
}

func TestFoundryDatasetPlanBodyPreservesComputedSHA256(t *testing.T) {
	customization := FoundryDatasetCustomization{}
	planBody, err := customization.PlanBodyFunc()(
		map[string]interface{}{"name": "example-dataset"},
		map[string]interface{}{"computed_sha256": "abc"},
	)
	if err != nil {
		t.Fatalf("PlanBodyFunc returned an error: %v", err)
	}
	values, ok := planBody.(map[string]interface{})
	if !ok || values["computed_sha256"] != "abc" {
		t.Fatalf("computed_sha256 was not preserved in plan body: %#v", planBody)
	}
}
