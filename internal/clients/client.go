package clients

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
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
	client.StopContext = ctx
	client.Features = o.Features

	azlog.SetListener(func(cls azlog.Event, msg string) {
		log.Printf("[DEBUG] %s %s: %s\n", time.Now().Format(time.StampMicro), cls, msg)
	})
	newResourceClient := NewNewResourceClient(o.SubscriptionId, o.Cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Telemetry: policy.TelemetryOptions{
				ApplicationID: o.ApplicationUserAgent,
			},
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
		},
		AuxiliaryTenants:      o.AuxiliaryTenantIDs,
		DisableRPRegistration: o.SkipProviderRegistration,
		Endpoint:              o.ARMEndpoint,
	})
	client.NewResourceClient = newResourceClient

	return nil
}
