package customization

import (
	"context"
	"net/http"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type KeyVaultKeyCustomization struct {
}

func (k KeyVaultKeyCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "create", id.ApiVersion, http.MethodPost, body, options)
		return err
	}
}

func (k KeyVaultKeyCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPatch, body, options)
		return err
	}
}

func (k KeyVaultKeyCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (k KeyVaultKeyCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.DeleteThenPoll(ctx, id, options)
		return err
	}
}

func (k KeyVaultKeyCustomization) GetResourceType() string {
	return "Microsoft.KeyVault/vaults/keys"
}

var _ DataPlaneResource = &KeyVaultKeyCustomization{}
