package customization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
		_, err = client.DataPlaneClient.Action(ctx, storageTableCollectionID(id), "", id.ApiVersion, http.MethodPost, payload, storageTableRequestOptions(options))
		return err
	}
}

func (c StorageTableCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, storageTableRequestOptions(options))
	}
}

func (c StorageTableCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		return fmt.Errorf("updating %q is not supported for %s; recreate the resource instead", id.Name, c.GetResourceType())
	}
}

func (c StorageTableCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.DeleteThenPoll(ctx, id, storageTableRequestOptions(options))
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
			payload[key] = value
		}
	}

	if rawName, ok := payload["TableName"]; ok {
		tableName, ok := rawName.(string)
		if !ok {
			return nil, fmt.Errorf(`body.TableName must be a string matching name %q`, id.Name)
		}
		if tableName != id.Name {
			return nil, fmt.Errorf(`body.TableName %q must match name %q`, tableName, id.Name)
		}
	}

	payload["TableName"] = id.Name
	return payload, nil
}

func storageTableRequestOptions(options clients.RequestOptions) clients.RequestOptions {
	cloned := options
	cloned.DisableAPIVersionQueryParameter = true
	cloned.APIVersionHeaderName = "x-ms-version"

	headers := make(map[string]string, len(options.Headers)+1)
	for key, value := range options.Headers {
		headers[key] = value
	}
	cloned.Headers = headers
	return cloned
}

var _ DataPlaneResource = &StorageTableCustomization{}
