package customization

import (
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func TestBuildStorageTableEntityBodyAddsCompositeKeys(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		Identifiers: map[string]string{
			"partitionKey": "pk",
			"rowKey":       "rk",
		},
	}

	body, err := buildStorageTableEntityBody(id, map[string]interface{}{
		"outputs": "value",
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got := body["PartitionKey"]; got != "pk" {
		t.Fatalf("expected PartitionKey to be injected, got %#v", got)
	}
	if got := body["RowKey"]; got != "rk" {
		t.Fatalf("expected RowKey to be injected, got %#v", got)
	}
}

func TestBuildStorageTableEntityBodyRejectsMismatchedKeys(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		Identifiers: map[string]string{
			"partitionKey": "pk",
			"rowKey":       "rk",
		},
	}

	_, err := buildStorageTableEntityBody(id, map[string]interface{}{
		"PartitionKey": "wrong",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), `identifiers.partitionKey "pk"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStorageTableEntityRequestOptionsSetTableHeaders(t *testing.T) {
	options := storageTableEntityRequestOptions(clients.RequestOptions{})
	if got := options.Headers["Accept"]; got != "application/json;odata=nometadata" {
		t.Fatalf("expected OData accept header, got %q", got)
	}
	if got := options.Headers["DataServiceVersion"]; got != "3.0;NetFx" {
		t.Fatalf("expected DataServiceVersion header, got %q", got)
	}
	if got := options.Headers["MaxDataServiceVersion"]; got != "3.0;NetFx" {
		t.Fatalf("expected MaxDataServiceVersion header, got %q", got)
	}
	if !options.DisableAPIVersionQueryParameter {
		t.Fatal("expected api-version query parameter to be disabled")
	}
}

func TestFlattenStorageTableEntityRemovesSystemFields(t *testing.T) {
	result, err := flattenStorageTableEntity(map[string]interface{}{
		"PartitionKey": "pk",
		"RowKey":       "rk",
		"Timestamp":    "ignored",
		"odata.etag":   "ignored",
		"outputs":      "value",
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	entity := result.(map[string]interface{})
	if _, ok := entity["PartitionKey"]; ok {
		t.Fatal("expected PartitionKey to be removed")
	}
	if _, ok := entity["RowKey"]; ok {
		t.Fatal("expected RowKey to be removed")
	}
	if got := entity["outputs"]; got != "value" {
		t.Fatalf("expected outputs to remain, got %#v", got)
	}
}
