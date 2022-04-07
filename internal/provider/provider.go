package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func AzureProvider() *schema.Provider {
	return azureProvider()
}

func azureProvider() *schema.Provider {
	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)

	resources["azapi_resource"] = services.ResourceAzApiResource()
	resources["azapi_update_resource"] = services.ResourceAzApiUpdateResource()

	dataSources["azapi_resource"] = services.ResourceAzApiDataSource()

	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
				Description: "The Subscription ID which should be used.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used.",
			},

			// TODO@mgd: this is blocked by https://github.com/Azure/azure-sdk-for-go/issues/17159
			// "auxiliary_tenant_ids": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	MaxItems: 3,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },

			"environment": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
				ValidateFunc: validation.StringInSlice([]string{"public", "usgovernment", "china"}, true),
				Description:  "The Cloud Environment which should be used. Possible values are public, usgovernment and china. Defaults to public.",
			},

			// TODO@mgd: the metadata_host is used to retrieve metadata from Azure to identify current environment, this is used to eliminate Azure Stack usage, in which case the provider doesn't support.
			// "metadata_host": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	DefaultFunc: schema.EnvDefaultFunc("ARM_METADATA_HOSTNAME", ""),
			// 	Description: "The Hostname which should be used for the Azure Metadata Service.",
			// },

			// Client Certificate specific fields
			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			// TODO@mgd: this depends on https://github.com/Azure/azure-sdk-for-go/pull/17099
			// "client_certificate_password": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
			// 	Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			// },

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"skip_provider_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_PROVIDER_REGISTRATION", false),
				Description: "Should the Provider skip registering all of the Resource Providers that it supports, if they're not already registered?",
			},

			// TODO@mgd: azidentity doesn't support msi_endpoint
			// // Managed Service Identity specific fields
			// "use_msi": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
			// 	Description: "Allowed Managed Service Identity be used for Authentication.",
			// },
			// "msi_endpoint": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
			// 	Description: "The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. ",
			// },

			"default_location": location.SchemaLocation(),

			"default_tags": tags.SchemaTags(),
		},

		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// var auxTenants []string
		// if v, ok := d.Get("auxiliary_tenant_ids").([]interface{}); ok && len(v) > 0 {
		// 	auxTenants = *utils.ExpandStringSlice(v)
		// } else if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
		// 	auxTenants = strings.Split(v, ";")
		// }

		var armEndpoint arm.Endpoint
		var authEndpoint azidentity.AuthorityHost
		env := d.Get("environment").(string)
		switch strings.ToLower(env) {
		case "public":
			armEndpoint = arm.AzurePublicCloud
			authEndpoint = azidentity.AzurePublicCloud
		case "usgovernment":
			armEndpoint = arm.AzureGovernment
			authEndpoint = azidentity.AzureGovernment
		case "china":
			armEndpoint = arm.AzureChina
			authEndpoint = azidentity.AzureChina
		default:
			return nil, diag.Errorf("unknown `environment` specified: %q", env)
		}

		// Maps the auth related environment variables used in the provider to what azidentity honors.
		os.Setenv("AZURE_TENANT_ID", d.Get("tenant_id").(string))
		os.Setenv("AZURE_CLIENT_ID", d.Get("client_id").(string))
		os.Setenv("AZURE_CLIENT_SECRET", d.Get("client_secret").(string))
		os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", d.Get("client_certificate_path").(string))

		cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
			AuthorityHost: authEndpoint,
			TenantID:      d.Get("tenant_id").(string),
		})
		if err != nil {
			return nil, diag.Errorf("failed to obtain a credential: %v", err)
		}

		copt := &clients.Option{
			SubscriptionId: d.Get("subscription_id").(string),
			Cred:           cred,
			//AuxiliaryTenantIDs:   auxTenants,
			ApplicationUserAgent: buildUserAgent(),
			ARMEndpoint:          armEndpoint,
			Features: features.UserFeatures{
				DefaultTags:     tags.ExpandTags(d.Get("default_tags").(map[string]interface{})),
				DefaultLocation: location.Normalize(d.Get("default_location").(string)),
			},
			SkipProviderRegistration: d.Get("skip_provider_registration").(bool),
		}

		//lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
		stopCtx, ok := schema.StopContext(ctx) //nolint:staticcheck
		if !ok {
			stopCtx = ctx
		}

		client := &clients.Client{}
		if err := client.Build(stopCtx, copt); err != nil {
			return nil, diag.FromErr(err)
		}

		// load schema
		var mutex sync.Mutex
		mutex.Lock()
		azure.GetAzureSchema()
		mutex.Unlock()
		return client, nil
	}
}

func buildUserAgent() string {
	// TODO: add more information to User Agent when length limitation is increased
	userAgent := "Terraform/azapi"

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		userAgent = fmt.Sprintf("%s/%s", userAgent, azureAgent)
	}
	return userAgent
}
