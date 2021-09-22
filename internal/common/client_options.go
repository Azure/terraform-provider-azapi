package common

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/sender"
)

type ClientOptions struct {
	SubscriptionId   string
	TenantID         string
	PartnerId        string
	TerraformVersion string

	ResourceManagerAuthorizer autorest.Authorizer
	ResourceManagerEndpoint   string
}

func (o ClientOptions) ConfigureClient(c *autorest.Client, authorizer autorest.Authorizer) {
	setUserAgent(c, o.TerraformVersion, o.PartnerId)

	c.Authorizer = authorizer
	c.Sender = sender.BuildSender("AzureRM")
}

func setUserAgent(client *autorest.Client, tfVersion, partnerID string) {
	// TODO:
}
