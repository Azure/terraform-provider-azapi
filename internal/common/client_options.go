package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/features"
	"github.com/ms-henglu/terraform-provider-azurermg/version"
)

type ClientOptions struct {
	SubscriptionId   string
	TenantID         string
	PartnerId        string
	TerraformVersion string

	ResourceManagerAuthorizer autorest.Authorizer
	ResourceManagerEndpoint   string

	Features features.UserFeatures
}

func (o ClientOptions) ConfigureClient(c *autorest.Client, authorizer autorest.Authorizer) {
	setUserAgent(c, o.TerraformVersion, o.PartnerId)

	c.Authorizer = authorizer
	c.Sender = sender.BuildSender("AzureRMG")
}

func setUserAgent(client *autorest.Client, tfVersion, partnerID string) {
	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", tfVersion, meta.SDKVersionString())

	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurermg/%s", tfUserAgent, version.ProviderVersion)
	client.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", client.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", client.UserAgent, azureAgent)
	}
}
