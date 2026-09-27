package customization

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

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
		// Use the Insert Entity operation (POST against the table collection) rather than
		// Insert-Or-Replace (PUT against the entity), so create fails with a conflict if the
		// entity already exists instead of silently overwriting it.
		_, err = client.DataPlaneClient.Action(ctx, storageTableEntityCollectionID(id), "", id.ApiVersion, http.MethodPost, payload, options)
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

// storageTableEntityCollectionID returns the table collection URL (without the composite-key
// suffix), used as the target for the Insert Entity operation.
func storageTableEntityCollectionID(id parse.DataPlaneResourceId) string {
	return storageTableEntityKeyPattern.ReplaceAllString(id.AzureResourceId, "")
}

func buildStorageTableEntityBody(id parse.DataPlaneResourceId, body interface{}) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if body != nil {
		bodyMap, ok := body.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected body for %s to be an object", id.AzureResourceType)
		}
		for key, value := range bodyMap {
			// PartitionKey and RowKey are derived entirely from parent_id, so reject them
			// outright rather than requiring them to match; strings.EqualFold catches case
			// variants such as "partitionkey" that the service would otherwise silently ignore.
			if strings.EqualFold(key, "PartitionKey") || strings.EqualFold(key, "RowKey") {
				return nil, fmt.Errorf(`body must not set %q; it is derived from parent_id`, key)
			}
			payload[key] = value
		}
	}

	partitionKey, rowKey, err := storageTableEntityKeys(id)
	if err != nil {
		return nil, err
	}

	payload["PartitionKey"] = partitionKey
	payload["RowKey"] = rowKey
	return payload, nil
}

var _ DataPlaneResource = &StorageTableEntityCustomization{}
