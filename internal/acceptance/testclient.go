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

		if v := os.Getenv("ARM_TENANT_ID"); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_TENANT_ID", v)
		}
		if v := os.Getenv("ARM_CLIENT_ID"); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_CLIENT_ID", v)
		}
		if v := os.Getenv("ARM_CLIENT_SECRET"); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_CLIENT_SECRET", v)
		}

		cred, err := azidentity.NewDefaultAzureCredential(
			&azidentity.DefaultAzureCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					Cloud: cloudConfig,
				},
			})
		if err != nil {
			return nil, fmt.Errorf("failed to obtain a credential: %v", err)
		}

		copt := &clients.Option{
			Cred:                     cred,
			CloudCfg:                 cloudConfig,
			Features:                 features.Default(),
			SkipProviderRegistration: true,
			TenantId:                 os.Getenv("ARM_TENANT_ID"),
			SubscriptionId:           os.Getenv("ARM_SUBSCRIPTION_ID"),
		}

		client := &clients.Client{}
		if err := client.Build(context.TODO(), copt); err != nil {
			return nil, err
		}
		_client = client
	}

	return _client, nil
}
