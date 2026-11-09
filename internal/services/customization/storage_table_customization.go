package customization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type StorageTableCustomization struct{}

func (c StorageTableCustomization) GetResourceType() string {
	return "Microsoft.Storage/storageAccounts/tableServices/tables"
}

func (c StorageTableCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		payload, err := buildStorageTableCreateBody(id, body)
		if err != nil {
			return err
		}

		ctx = storageTablePollingContext(ctx, options.Headers)

		_, err = client.DataPlaneClient.Action(ctx, storageTableCollectionID(id), "", id.ApiVersion, http.MethodPost, payload, options)
		return err
	}
}

func (c StorageTableCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (c StorageTableCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		return fmt.Errorf("updating %q is not supported for %s; recreate the resource instead", id.Name, c.GetResourceType())
	}
}

func (c StorageTableCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.DeleteThenPoll(ctx, id, options)
		return err
	}
}

func storageTableCollectionID(id parse.DataPlaneResourceId) string {
	return strings.TrimSuffix(id.ParentId, "/") + "/Tables"
}

func buildStorageTableCreateBody(id parse.DataPlaneResourceId, body interface{}) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if body != nil {
		bodyMap, ok := body.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected body for %s to be an object", id.AzureResourceType)
		}
		for key, value := range bodyMap {
			// TableName comes from name; EqualFold also rejects case variants such as "tableName".
			if strings.EqualFold(key, "TableName") {
				return nil, fmt.Errorf(`body must not set %q; it is derived from name %q`, key, id.Name)
			}
			payload[key] = value
		}
	}

	payload["TableName"] = id.Name
	return payload, nil
}

// storageTablePollingContext propagates x-ms-version and Accept to SDK-generated polling requests.
// Without them, Table Storage rejects the polling GET with 403 (missing version) or 415 (Atom format not supported).
func storageTablePollingContext(ctx context.Context, headers map[string]string) context.Context {
	pollingHeaders := make(http.Header)
	pollingHeaders.Set("Accept", "application/json")
	for key, value := range headers {
		if strings.EqualFold(key, "x-ms-version") || strings.EqualFold(key, "Accept") {
			pollingHeaders.Set(key, value)
		}
	}
	return policy.WithHTTPHeader(ctx, pollingHeaders)
}

var _ DataPlaneResource = &StorageTableCustomization{}
