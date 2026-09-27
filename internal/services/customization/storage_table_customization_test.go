package customization

import (
	"strings"
	"testing"

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

func TestBuildStorageTableCreateBodyRejectsTableNameInput(t *testing.T) {
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
	if !strings.Contains(err.Error(), `must not set "TableName"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildStorageTableCreateBodyRejectsTableNameInputCaseInsensitively(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables",
		Name:              "acctesttable",
	}

	_, err := buildStorageTableCreateBody(id, map[string]interface{}{
		"tablename": "acctesttable",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), `must not set "tablename"`) {
		t.Fatalf("unexpected error: %v", err)
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
