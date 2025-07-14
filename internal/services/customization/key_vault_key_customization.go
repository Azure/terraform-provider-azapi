package customization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
)

type KeyVaultKeyCustomization struct {
}

func (k KeyVaultKeyCustomization) GetResourceType() string {
	return "Microsoft.KeyVault/vaults/keys"
}

func (k KeyVaultKeyCustomization) CreateFunc() Func {
	return nil
}

func (k KeyVaultKeyCustomization) UpdateFunc() Func {
	return nil
}

func (k KeyVaultKeyCustomization) ReadFunc() Func {
	return nil
}

func (k KeyVaultKeyCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, clients clients.Client, id parse.ResourceId, options clients.RequestOptions) error {

		dataPlaneClient := clients.DataPlaneClient

		path := id.AzureResourceId
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimSuffix(path, "/")
		components := strings.Split(path, "/")
		parts := make(map[string]string)
		for i := 0; i < len(components)-1; i += 2 {
			parts[components[i]] = components[i+1]
		}

		if parts["vaults"] == "" {
			return fmt.Errorf("key vault name is missing in the resource ID: %s", id.AzureResourceId)
		}

		resourceID := fmt.Sprintf("%s.vault.azure.net/keys/%s", parts["vaults"], id.Name)

		_, err := dataPlaneClient.Action(ctx, resourceID, "", "7.4", http.MethodDelete, nil, options)
		if err != nil && !utils.ResponseErrorWasNotFound(err) {
			return err
		}
		return nil
	}
}

var _ Resource = &KeyVaultKeyCustomization{}
