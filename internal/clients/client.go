package clients

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/features"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Features features.UserFeatures

	NewResourceClient *NewResourceClient
}

type Option struct {
	SubscriptionId           string
	Cred                     azcore.TokenCredential
	ARMEndpoint              arm.Endpoint
	AuxiliaryTenantIDs       []string
	ApplicationUserAgent     string
	Features                 features.UserFeatures
	SkipProviderRegistration bool
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *Option) error {
	// TODO@mgd: whether we need to regard 429 as non-retriable error in this provider?
	//autorest.Count429AsRetry = false
	client.StopContext = ctx
	client.Features = o.Features

	newResourceClient := NewNewResourceClient(o.SubscriptionId, o.Cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Telemetry: policy.TelemetryOptions{
				ApplicationID: o.ApplicationUserAgent,
			},
		},
		AuxiliaryTenants:      o.AuxiliaryTenantIDs,
		DisableRPRegistration: o.SkipProviderRegistration,
		Endpoint:              o.ARMEndpoint,
	})
	client.NewResourceClient = newResourceClient

	return nil
}
