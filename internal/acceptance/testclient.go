package acceptance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/vcr"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
)

var (
	_client    *clients.Client
	clientLock = &sync.Mutex{}
)

type TestClientOptions struct {
	Transport policy.Transporter

	Credential azcore.TokenCredential

	SubscriptionID string
	TenantID       string
}

func BuildTestClient() (*clients.Client, error) {
	if vcr.Enabled() {
		if rec := vcr.Active(); rec != nil {
			return BuildTestClientWithOptions(TestClientOptions{
				Transport:      rec.Transport(),
				Credential:     rec.Credential(),
				SubscriptionID: rec.SubscriptionID(),
				TenantID:       rec.TenantID(),
			})
		}
	}

	clientLock.Lock()
	defer clientLock.Unlock()

	if _client == nil {
		client, err := BuildTestClientWithOptions(TestClientOptions{})
		if err != nil {
			return nil, err
		}
		_client = client
	}

	return _client, nil
}

func BuildTestClientWithOptions(o TestClientOptions) (*clients.Client, error) {
	cloudConfig, err := testCloudConfig()
	if err != nil {
		return nil, err
	}

	cred := o.Credential
	if cred == nil {
		cred, err = buildTestCredential(cloudConfig)
		if err != nil {
			return nil, err
		}
	}

	subscriptionID := o.SubscriptionID
	if subscriptionID == "" {
		subscriptionID = os.Getenv("ARM_SUBSCRIPTION_ID")
	}
	tenantID := o.TenantID
	if tenantID == "" {
		tenantID = os.Getenv("ARM_TENANT_ID")
	}

	copt := &clients.Option{
		Cred:                     cred,
		CloudCfg:                 cloudConfig,
		Features:                 features.Default(),
		SkipProviderRegistration: true,
		TenantId:                 tenantID,
		SubscriptionId:           subscriptionID,
		Transport:                o.Transport,
	}

	client := &clients.Client{}
	if err := client.Build(context.TODO(), copt); err != nil {
		return nil, err
	}

	return client, nil
}

func testCloudConfig() (cloud.Configuration, error) {
	env := os.Getenv("ARM_ENVIRONMENT")
	switch strings.ToLower(env) {
	case "public":
		return cloud.AzurePublic, nil
	case "usgovernment":
		return cloud.AzureGovernment, nil
	case "china":
		return cloud.AzureChina, nil
	case "custom":
		activeDirectoryAuthorityHost := os.Getenv("ARM_ACTIVE_DIRECTORY_AUTHORITY_HOST")
		if activeDirectoryAuthorityHost == "" {
			return cloud.Configuration{}, fmt.Errorf("`ARM_ACTIVE_DIRECTORY_AUTHORITY_HOST` must be set when `environment` is set to `custom`")
		}
		resourceManagerEndpoint := os.Getenv("ARM_RESOURCE_MANAGER_ENDPOINT")
		if resourceManagerEndpoint == "" {
			return cloud.Configuration{}, fmt.Errorf("`ARM_RESOURCE_MANAGER_ENDPOINT` must be set when `environment` is set to `custom`")
		}
		resourceManagerAudience := os.Getenv("ARM_RESOURCE_MANAGER_AUDIENCE")
		if resourceManagerAudience == "" {
			return cloud.Configuration{}, fmt.Errorf("`ARM_RESOURCE_MANAGER_AUDIENCE` must be set when `environment` is set to `custom`")
		}

		return cloud.Configuration{
			ActiveDirectoryAuthorityHost: activeDirectoryAuthorityHost,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: resourceManagerAudience,
					Endpoint: resourceManagerEndpoint,
				},
			},
		}, nil
	default:
		return cloud.AzurePublic, nil
	}
}

// buildTestCredential mirrors BuildTestClient's default credential behavior: it
// propagates ARM_* values to the AZURE_* variables expected by azidentity and
// returns a DefaultAzureCredential.
func buildTestCredential(cloudConfig cloud.Configuration) (azcore.TokenCredential, error) {
	if v := os.Getenv("ARM_TENANT_ID"); len(v) != 0 {
		if err := os.Setenv("AZURE_TENANT_ID", v); err != nil {
			return nil, fmt.Errorf("setting AZURE_TENANT_ID: %w", err)
		}
	}
	if v := os.Getenv("ARM_CLIENT_ID"); len(v) != 0 {
		if err := os.Setenv("AZURE_CLIENT_ID", v); err != nil {
			return nil, fmt.Errorf("setting AZURE_CLIENT_ID: %w", err)
		}
	}
	if v := os.Getenv("ARM_CLIENT_SECRET"); len(v) != 0 {
		if err := os.Setenv("AZURE_CLIENT_SECRET", v); err != nil {
			return nil, fmt.Errorf("setting AZURE_CLIENT_SECRET: %w", err)
		}
	}

	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: cloudConfig,
			},
			DisableInstanceDiscovery: strings.EqualFold(os.Getenv("ARM_DISABLE_INSTANCE_DISCOVERY"), "true"),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
	}

	return cred, nil
}
