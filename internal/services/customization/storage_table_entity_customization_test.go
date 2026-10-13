package customization

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
)

func TestBuildStorageTableEntityBodyAddsCompositeKeys(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		AzureResourceId:   "mystorage.table.core.windows.net/mytable(PartitionKey='pk',RowKey='rk')",
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

func TestBuildStorageTableEntityBodyRejectsPartitionKeyInput(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		AzureResourceId:   "mystorage.table.core.windows.net/mytable(PartitionKey='pk',RowKey='rk')",
	}

	_, err := buildStorageTableEntityBody(id, map[string]interface{}{
		"PartitionKey": "pk",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), `must not set "PartitionKey"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildStorageTableEntityBodyRejectsRowKeyInputCaseInsensitively(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		AzureResourceId:   "mystorage.table.core.windows.net/mytable(PartitionKey='pk',RowKey='rk')",
	}

	_, err := buildStorageTableEntityBody(id, map[string]interface{}{
		"rowkey": "rk",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), `must not set "rowkey"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildStorageTableEntityBodyRejectsMissingKeys(t *testing.T) {
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		AzureResourceId:   "mystorage.table.core.windows.net/mytable",
		ParentId:          "mystorage.table.core.windows.net/mytable",
	}

	_, err := buildStorageTableEntityBody(id, map[string]interface{}{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "must end with (PartitionKey=") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStorageTableEntityCreateInsertsNewEntity(t *testing.T) {
	transport := &storageTableEntityTransport{}
	dataPlaneClient, err := clients.NewDataPlaneClient(storageTableStaticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("building data plane client: %v", err)
	}

	customization := StorageTableEntityCustomization{}
	client := clients.Client{
		DataPlaneClient: dataPlaneClient,
	}
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		AzureResourceId:   "management.azure.com/mytable(PartitionKey='pk',RowKey='rk')",
		ApiVersion:        "2026-04-06",
	}
	options := clients.RequestOptions{
		Headers: map[string]string{
			"x-ms-version": "2026-04-06",
		},
	}
	createBody := map[string]interface{}{
		"retained": "initial",
	}
	if err := customization.CreateFunc()(context.Background(), client, id, createBody, options); err != nil {
		t.Fatalf("creating entity: %v", err)
	}

	if len(transport.requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(transport.requests))
	}
	got := transport.requests[0]
	if got.method != http.MethodPost {
		t.Fatalf("expected create to use the Insert Entity POST operation, got %q", got.method)
	}
	if got.url != "https://management.azure.com/mytable" {
		t.Fatalf("expected Insert Entity request against the table collection, got %q", got.url)
	}
	expectedPayload := map[string]interface{}{
		"PartitionKey": "pk",
		"RowKey":       "rk",
		"retained":     "initial",
	}
	if !reflect.DeepEqual(got.body, expectedPayload) {
		t.Fatalf("expected request body %#v, got %#v", expectedPayload, got.body)
	}
}

func TestStorageTableEntityUpdatePutReplacesEntity(t *testing.T) {
	transport := &storageTableEntityTransport{
		entity: map[string]interface{}{
			"retained":         "initial",
			"removed_property": "remove-me",
		},
	}
	dataPlaneClient, err := clients.NewDataPlaneClient(storageTableStaticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("building data plane client: %v", err)
	}

	customization := StorageTableEntityCustomization{}
	client := clients.Client{
		DataPlaneClient: dataPlaneClient,
	}
	id := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
		AzureResourceId:   "management.azure.com/mytable(PartitionKey='pk',RowKey='rk')",
		ApiVersion:        "2026-04-06",
	}
	options := clients.RequestOptions{
		Headers: map[string]string{
			"x-ms-version": "2026-04-06",
			"x-custom":     "custom-value",
		},
	}

	transport.entity["external_property"] = "remove-on-next-write"
	updateBody := map[string]interface{}{
		"retained": "updated",
		"added":    "new",
	}
	if err := customization.UpdateFunc()(context.Background(), client, id, updateBody, options); err != nil {
		t.Fatalf("updating entity: %v", err)
	}

	if len(transport.requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(transport.requests))
	}
	got := transport.requests[0]
	if got.method != http.MethodPut {
		t.Fatalf("expected update to use the Insert-Or-Replace PUT operation, got %q", got.method)
	}
	if got.url != "https://management.azure.com/mytable(PartitionKey='pk',RowKey='rk')" {
		t.Fatalf("expected Insert-Or-Replace request against the entity, got %q", got.url)
	}
	expectedUpdatePayload := map[string]interface{}{
		"PartitionKey": "pk",
		"RowKey":       "rk",
		"retained":     "updated",
		"added":        "new",
	}
	if !reflect.DeepEqual(got.body, expectedUpdatePayload) {
		t.Fatalf("expected request body %#v, got %#v", expectedUpdatePayload, got.body)
	}
	if gotHeader := got.headers.Get("x-ms-version"); gotHeader != "2026-04-06" {
		t.Fatalf("expected caller-provided x-ms-version header, got %q", gotHeader)
	}
	if gotHeader := got.headers.Get("x-custom"); gotHeader != "custom-value" {
		t.Fatalf("expected caller-provided custom header, got %q", gotHeader)
	}
	if gotHeader := got.headers.Get("If-Match"); gotHeader != "" {
		t.Fatalf("expected no If-Match header, got %q", gotHeader)
	}

	result, err := customization.ReadFunc()(context.Background(), client, id, options)
	if err != nil {
		t.Fatalf("reading entity: %v", err)
	}
	expected := map[string]interface{}{
		"PartitionKey": "pk",
		"RowKey":       "rk",
		"Timestamp":    "2026-09-13T00:00:00Z",
		"odata.etag":   `W/"datetime'2026-09-13T00%3A00%3A00.0000000Z'"`,
		"retained":     "updated",
		"added":        "new",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected full response %#v, got %#v", expected, result)
	}
	for _, removed := range []string{"removed_property", "external_property"} {
		if _, ok := result.(map[string]interface{})[removed]; ok {
			t.Fatalf("expected %q to be removed by replacement", removed)
		}
	}

	managedBody := utils.UpdateObject(updateBody, result, utils.UpdateJsonOption{IgnoreMissingProperty: true})
	if !reflect.DeepEqual(managedBody, updateBody) {
		t.Fatalf("expected managed body to remain configuration-shaped, got %#v", managedBody)
	}
	refreshedBody := utils.UpdateObject(managedBody, result, utils.UpdateJsonOption{IgnoreMissingProperty: true})
	if !reflect.DeepEqual(refreshedBody, managedBody) {
		t.Fatalf("expected repeated refresh to produce no managed-body diff, got %#v", refreshedBody)
	}
}

type storageTableStaticTokenCredential struct{}

func (storageTableStaticTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "test-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

type storageTableEntityRequest struct {
	method  string
	url     string
	headers http.Header
	body    map[string]interface{}
}

type storageTableEntityTransport struct {
	requests []storageTableEntityRequest
	entity   map[string]interface{}
}

func (t *storageTableEntityTransport) Do(request *http.Request) (*http.Response, error) {
	switch request.Method {
	case http.MethodPost, http.MethodPut:
		var body map[string]interface{}
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decoding request body: %w", err)
		}
		t.requests = append(t.requests, storageTableEntityRequest{
			method:  request.Method,
			url:     fmt.Sprintf("%s://%s%s", request.URL.Scheme, request.URL.Host, request.URL.Path),
			headers: request.Header.Clone(),
			body:    body,
		})
		t.entity = body
		// Insert Entity (POST) responds 201 Created; Insert-Or-Replace Entity (PUT) responds
		// 204 No Content.
		statusCode := http.StatusNoContent
		if request.Method == http.MethodPost {
			statusCode = http.StatusCreated
		}
		return storageTableResponse(request, statusCode, "")
	case http.MethodGet:
		responseBody := make(map[string]interface{}, len(t.entity)+2)
		for key, value := range t.entity {
			responseBody[key] = value
		}
		responseBody["Timestamp"] = "2026-09-13T00:00:00Z"
		responseBody["odata.etag"] = `W/"datetime'2026-09-13T00%3A00%3A00.0000000Z'"`
		body, err := json.Marshal(responseBody)
		if err != nil {
			return nil, fmt.Errorf("encoding response body: %w", err)
		}
		return storageTableResponse(request, http.StatusOK, string(body))
	default:
		return storageTableResponse(request, http.StatusMethodNotAllowed, "")
	}
}

func storageTableResponse(request *http.Request, statusCode int, body string) (*http.Response, error) {
	return &http.Response{
		StatusCode: statusCode,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(body)),
		Request: request,
	}, nil
}
