package customization

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type StorageTableEntityCustomization struct{}

func (c StorageTableEntityCustomization) GetResourceType() string {
	return "Microsoft.Storage/storageAccounts/tableServices/tables/entities"
}

func (c StorageTableEntityCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		payload, err := buildStorageTableEntityBody(id, body)
		if err != nil {
			return err
		}
		_, err = client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPut, payload, options)
		return err
	}
}

func (c StorageTableEntityCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (c StorageTableEntityCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		payload, err := buildStorageTableEntityBody(id, body)
		if err != nil {
			return err
		}
		_, err = client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPut, payload, options)
		return err
	}
}

func (c StorageTableEntityCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.DeleteThenPoll(ctx, id, options)
		return err
	}
}

type StorageTableEntitiesCustomization struct{}

func (c StorageTableEntitiesCustomization) GetResourceType() string {
	return "Microsoft.Storage/storageAccounts/tableServices/tables/entitiesCollection"
}

func (c StorageTableEntitiesCustomization) CreateFunc() CreateFunc { return nil }
func (c StorageTableEntitiesCustomization) UpdateFunc() UpdateFunc { return nil }
func (c StorageTableEntitiesCustomization) DeleteFunc() DeleteFunc { return nil }

func (c StorageTableEntitiesCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

// storageTableEntityKeyPattern matches the composite key suffix embedded in the entity's
// parent_id / resource ID, e.g. "mytable(PartitionKey='Sales',RowKey='1')".
var storageTableEntityKeyPattern = regexp.MustCompile(`\(PartitionKey='([^']*)',RowKey='([^']*)'\)$`)

// storageTableEntityKeys extracts PartitionKey and RowKey from the entity resource ID.
func storageTableEntityKeys(id parse.DataPlaneResourceId) (string, string, error) {
	match := storageTableEntityKeyPattern.FindStringSubmatch(id.AzureResourceId)
	if match == nil {
		return "", "", fmt.Errorf("parent_id for %s must end with (PartitionKey='<value>',RowKey='<value>'), got %q", id.AzureResourceType, id.ParentId)
	}
	return match[1], match[2], nil
}

func buildStorageTableEntityBody(id parse.DataPlaneResourceId, body interface{}) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if body != nil {
		bodyMap, ok := body.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected body for %s to be an object", id.AzureResourceType)
		}
		for key, value := range bodyMap {
			payload[key] = value
		}
	}

	partitionKey, rowKey, err := storageTableEntityKeys(id)
	if err != nil {
		return nil, err
	}

	if rawPartitionKey, ok := payload["PartitionKey"]; ok {
		value, ok := rawPartitionKey.(string)
		if !ok || value != partitionKey {
			return nil, fmt.Errorf(`body.PartitionKey must be a string matching the PartitionKey %q in parent_id`, partitionKey)
		}
	}
	if rawRowKey, ok := payload["RowKey"]; ok {
		value, ok := rawRowKey.(string)
		if !ok || value != rowKey {
			return nil, fmt.Errorf(`body.RowKey must be a string matching the RowKey %q in parent_id`, rowKey)
		}
	}

	payload["PartitionKey"] = partitionKey
	payload["RowKey"] = rowKey
	return payload, nil
}

var _ DataPlaneResource = &StorageTableEntityCustomization{}
var _ DataPlaneResource = &StorageTableEntitiesCustomization{}
