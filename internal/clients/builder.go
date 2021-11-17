package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/common"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/features"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
)

type ClientBuilder struct {
	AuthConfig                  *authentication.Config
	DisableCorrelationRequestID bool
	CustomCorrelationRequestID  string
	DisableTerraformPartnerID   bool
	PartnerId                   string
	SkipProviderRegistration    bool
	StorageUseAzureAD           bool
	TerraformVersion            string
	Features                    features.UserFeatures
}

const azureStackEnvironmentError = `
The AzureRM Provider supports the different Azure Public Clouds - including China, Germany,
Public and US Government - however it does not support Azure Stack due to differences in
API and feature availability.

Terraform instead offers a separate "azurestack" provider which supports the functionality
and API's available in Azure Stack via Azure Stack Profiles.
`

func Build(ctx context.Context, builder ClientBuilder) (*Client, error) {
	// point folks towards the separate Azure Stack Provider when using Azure Stack
	if strings.EqualFold(builder.AuthConfig.Environment, "AZURESTACKCLOUD") {
		return nil, fmt.Errorf(azureStackEnvironmentError)
	}

	isAzureStack, err := authentication.IsEnvironmentAzureStack(ctx, builder.AuthConfig.MetadataHost, builder.AuthConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to determine if environment is Azure Stack: %+v", err)
	}
	if isAzureStack {
		return nil, fmt.Errorf(azureStackEnvironmentError)
	}

	env, err := authentication.AzureEnvironmentByNameFromEndpoint(ctx, builder.AuthConfig.MetadataHost, builder.AuthConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to find environment %q from endpoint %q: %+v", builder.AuthConfig.Environment, builder.AuthConfig.MetadataHost, err)
	}

	client := Client{}

	oauthConfig, err := builder.AuthConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("building OAuth Config: %+v", err)
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("unable to configure OAuthConfig for tenant %s", builder.AuthConfig.TenantID)
	}

	sender := sender.BuildSender("AzureRM")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, fmt.Errorf("unable to get authorization token for resource manager: %+v", err)
	}

	o := &common.ClientOptions{
		SubscriptionId:            builder.AuthConfig.SubscriptionID,
		TenantID:                  builder.AuthConfig.TenantID,
		PartnerId:                 builder.PartnerId,
		TerraformVersion:          builder.TerraformVersion,
		ResourceManagerAuthorizer: auth,
		ResourceManagerEndpoint:   endpoint,
		Features:                  builder.Features,
	}

	if err := client.Build(ctx, o); err != nil {
		return nil, fmt.Errorf("building Client: %+v", err)
	}

	return &client, nil
}
