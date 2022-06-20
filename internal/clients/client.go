package clients

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/features"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Features features.UserFeatures

	ResourceClient *ResourceClient
}

type Option struct {
	SubscriptionId              string
	Cred                        azcore.TokenCredential
	ApplicationUserAgent        string
	Features                    features.UserFeatures
	SkipProviderRegistration    bool
	DisableCorrelationRequestID bool
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *Option) error {
	client.StopContext = ctx
	client.Features = o.Features

	azlog.SetListener(func(cls azlog.Event, msg string) {
		log.Printf("[DEBUG] %s %s: %s\n", time.Now().Format(time.StampMicro), cls, msg)
	})

	perCallPolicies := make([]policy.Policy, 0)
	perCallPolicies = append(perCallPolicies, withUserAgent(o.ApplicationUserAgent))
	if !o.DisableCorrelationRequestID {
		perCallPolicies = append(perCallPolicies, withCorrelationRequestID(correlationRequestID()))
	}

	resourceClient, err := NewResourceClient(o.SubscriptionId, o.Cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			// Disable the default telemetry policy, because it has a length limitation for user agent
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
			PerCallPolicies: perCallPolicies,
		},
		DisableRPRegistration: o.SkipProviderRegistration,
	})
	if err != nil {
		return err
	}
	client.ResourceClient = resourceClient

	return nil
}
