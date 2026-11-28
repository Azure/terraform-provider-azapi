package customization

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
)

// KeyVaultKeyRotationPolicyCustomization supports the rotation policy of a Key Vault key.
//
// Create: PUT {vault}/keys/{name}/rotationpolicy
// Get:    GET {vault}/keys/{name}/rotationpolicy
// Update: PUT {vault}/keys/{name}/rotationpolicy
// Delete: PUT {vault}/keys/{name}/rotationpolicy with no lifetime actions, as the AzureRM provider does
//
// Every key has a rotation policy, so create doesn't check whether it exists. There is no delete
// operation, and a DELETE request to {vault}/keys/{name}/rotationpolicy deletes the key itself.
type KeyVaultKeyRotationPolicyCustomization struct{}

func (k KeyVaultKeyRotationPolicyCustomization) GetResourceType() string {
	return "Microsoft.KeyVault/vaults/keys/rotationPolicy"
}

func (k KeyVaultKeyRotationPolicyCustomization) CreateFunc() CreateFunc {
	return nil
}

func (k KeyVaultKeyRotationPolicyCustomization) ReadFunc() ReadFunc {
	return nil
}

func (k KeyVaultKeyRotationPolicyCustomization) UpdateFunc() UpdateFunc {
	return nil
}

func (k KeyVaultKeyRotationPolicyCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		policyWithoutActions := map[string]interface{}{
			"lifetimeActions": []interface{}{},
		}
		_, err := client.DataPlaneClient.CreateOrUpdateThenPoll(ctx, id, policyWithoutActions, options)
		if err != nil && utils.ResponseErrorWasNotFound(err) {
			// The key has been deleted, together with its rotation policy.
			return nil
		}
		return err
	}
}

func (k KeyVaultKeyRotationPolicyCustomization) SkipExistenceCheck() bool {
	return true
}

var _ DataPlaneResource = &KeyVaultKeyRotationPolicyCustomization{}
var _ DataPlaneResourceWithoutExistenceCheck = &KeyVaultKeyRotationPolicyCustomization{}
