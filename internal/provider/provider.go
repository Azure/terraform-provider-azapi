package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/clients"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/features"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AzureProvider() *schema.Provider {
	return azureProvider()
}

func azureProvider() *schema.Provider {
	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)

	resources["azurerm-restapi_resource"] = services.ResourceAzureGenericResource()
	resources["azurerm-restapi_patch_resource"] = services.ResourceAzureGenericPatchResource()

	dataSources["azurerm-restapi_resource"] = services.ResourceAzureGenericDataSource()

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

			"auxiliary_tenant_ids": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
				Description: "The Cloud Environment which should be used. Possible values are public, usgovernment, german, and china. Defaults to public.",
			},

			"metadata_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_METADATA_HOSTNAME", ""),
				Description: "The Hostname which should be used for the Azure Metadata Service.",
			},

			// Client Certificate specific fields
			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
				Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			// Managed Service Identity specific fields
			"use_msi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
				Description: "Allowed Managed Service Identity be used for Authentication.",
			},
			"msi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
				Description: "The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. ",
			},

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
		var auxTenants []string
		if v, ok := d.Get("auxiliary_tenant_ids").([]interface{}); ok && len(v) > 0 {
			auxTenants = *utils.ExpandStringSlice(v)
		} else if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
			auxTenants = strings.Split(v, ";")
		}

		metadataHost := d.Get("metadata_host").(string)

		builder := &authentication.Builder{
			SubscriptionID:     d.Get("subscription_id").(string),
			ClientID:           d.Get("client_id").(string),
			ClientSecret:       d.Get("client_secret").(string),
			TenantID:           d.Get("tenant_id").(string),
			AuxiliaryTenantIDs: auxTenants,
			Environment:        d.Get("environment").(string),
			MetadataHost:       metadataHost,
			MsiEndpoint:        d.Get("msi_endpoint").(string),
			ClientCertPassword: d.Get("client_certificate_password").(string),
			ClientCertPath:     d.Get("client_certificate_path").(string),

			// Feature Toggles
			SupportsClientCertAuth:         true,
			SupportsClientSecretAuth:       true,
			SupportsManagedServiceIdentity: d.Get("use_msi").(bool),
			SupportsAzureCliToken:          true,

			// Doc Links
			ClientSecretDocsLink: "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret",
		}

		config, err := builder.Build()
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("building AzureRM Client: %s", err))
		}

		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		clientBuilder := clients.ClientBuilder{
			AuthConfig:       config,
			TerraformVersion: terraformVersion,
			Features: features.UserFeatures{
				DefaultTags:     tags.ExpandTags(d.Get("default_tags").(map[string]interface{})),
				DefaultLocation: location.Normalize(d.Get("default_location").(string)),
			},
		}

		//lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
		stopCtx, ok := schema.StopContext(ctx) //nolint:staticcheck
		if !ok {
			stopCtx = ctx
		}

		client, err := clients.Build(stopCtx, clientBuilder)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client.StopContext = stopCtx

		// load schema
		var mutex sync.Mutex
		mutex.Lock()
		azure.GetAzureSchema()
		mutex.Unlock()
		return client, nil
	}
}
