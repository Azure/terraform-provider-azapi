package customization

import (
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func TestBuildStorageTableCreateBodyAddsTableName(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables",
		Name:              "acctesttable",
	}

	body, err := buildStorageTableCreateBody(id, map[string]interface{}{})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got := body["TableName"]; got != "acctesttable" {
		t.Fatalf("expected TableName to be set from name, got: %#v", got)
	}
}

func TestBuildStorageTableCreateBodyRejectsMismatchedTableName(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables",
		Name:              "acctesttable",
	}

	_, err := buildStorageTableCreateBody(id, map[string]interface{}{
		"TableName": "different",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), `must match name "acctesttable"`) {
		t.Fatalf("expected name mismatch error, got: %v", err)
	}
}

func TestStorageTableRequestOptionsUseVersionHeader(t *testing.T) {
	options := storageTableRequestOptions(clients.RequestOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})

	if !options.DisableAPIVersionQueryParameter {
		t.Fatal("expected api-version query parameter to be disabled")
	}
	if options.APIVersionHeaderName != "x-ms-version" {
		t.Fatalf("expected x-ms-version header, got %q", options.APIVersionHeaderName)
	}
	if got := options.Headers["Accept"]; got != "application/json" {
		t.Fatalf("expected existing header to be preserved, got %q", got)
	}
}

func TestGetCustomizationStorageTable(t *testing.T) {
	customization := GetCustomization("Microsoft.Storage/storageAccounts/tableServices/tables@2025-11-05")
	if customization == nil {
		t.Fatal("expected storage table customization to be registered")
	}
	if got := (*customization).GetResourceType(); got != "Microsoft.Storage/storageAccounts/tableServices/tables" {
		t.Fatalf("unexpected resource type %q", got)
	}
}
