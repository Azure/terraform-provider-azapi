package clients

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/common"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/features"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Features features.UserFeatures

	ResourceClient *ResourceClient
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true
	client.StopContext = ctx

	client.Features = o.Features

	resourceClient := NewResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&resourceClient.Client, o.ResourceManagerAuthorizer)
	client.ResourceClient = &resourceClient

	return nil
}
