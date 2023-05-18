package acceptance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
)

var (
	_client    *clients.Client
	clientLock = &sync.Mutex{}
)

func BuildTestClient() (*clients.Client, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if _client == nil {
		var cloudConfig cloud.Configuration
		env := os.Getenv("ARM_ENVIRONMENT")
		switch strings.ToLower(env) {
		case "public":
			cloudConfig = cloud.AzurePublic
		case "usgovernment":
			cloudConfig = cloud.AzureGovernment
		case "china":
			cloudConfig = cloud.AzureChina
		default:
			cloudConfig = cloud.AzurePublic
		}

		cred, err := azidentity.NewClientSecretCredential(
			os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"),
			&azidentity.ClientSecretCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					Cloud: cloudConfig,
				},
			})
		if err != nil {
			return nil, fmt.Errorf("failed to obtain a credential: %v", err)
		}

		copt := &clients.Option{
			Cred:                     cred,
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
