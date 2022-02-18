package acceptance

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/clients"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/features"
)

var (
	_client    *clients.Client
	clientLock = &sync.Mutex{}
)

func BuildTestClient() (*clients.Client, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if _client == nil {
		var armEndpoint arm.Endpoint
		var authEndpoint azidentity.AuthorityHost
		env := os.Getenv("ARM_ENVIRONMENT")
		switch env {
		case "public":
			armEndpoint = arm.AzurePublicCloud
			authEndpoint = azidentity.AzurePublicCloud
		case "usgovernment":
			armEndpoint = arm.AzureGovernment
			authEndpoint = azidentity.AzureGovernment
		case "china":
			armEndpoint = arm.AzureChina
			authEndpoint = azidentity.AzureChina
		default:
			armEndpoint = arm.AzurePublicCloud
			authEndpoint = azidentity.AzurePublicCloud
		}

		cred, err := azidentity.NewClientSecretCredential(
			os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"),
			&azidentity.ClientSecretCredentialOptions{
				AuthorityHost: authEndpoint,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to obtain a credential: %v", err)
		}

		copt := &clients.Option{
			SubscriptionId:           os.Getenv("ARM_SUBSCRIPTION_ID"),
			Cred:                     cred,
			ARMEndpoint:              armEndpoint,
			Features:                 features.Default(),
			SkipProviderRegistration: true,
		}

		client := &clients.Client{}
		if err := client.Build(context.TODO(), copt); err != nil {
			return nil, err
		}
		_client = client
	}

	return _client, nil
}
