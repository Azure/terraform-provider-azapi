package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/services"

	"github.com/Azure/terraform-provider-azapi/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
)

func AzureProvider() *schema.Provider {
	return azureProvider()
}

func azureProvider() *schema.Provider {
	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)

	resources["azapi_resource"] = services.ResourceAzApiResource()
	resources["azapi_update_resource"] = services.ResourceAzApiUpdateResource()
	resources["azapi_resource_action"] = services.ResourceResourceAction()

	dataSources["azapi_resource"] = services.ResourceAzApiDataSource()
	dataSources["azapi_resource_action"] = services.ResourceResourceActionDataSource()

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

			// OIDC specifc fields
			"oidc_token_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", ""),
				Description: "The file where the token from the OIDC provider is. For use When authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_authority_host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_AUTHORITY_HOST", ""),
				Description: "The authority host for OIDC login. For use When authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, ""),
				Description: "The token from the OIDC provider. For use When authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, ""),
				Description: "The bearer token for the request to the OIDC provider. For use When authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"}, ""),
				Description: "The URL for the OIDC provider from which to request an ID token. For use When authenticating as a Service Principal using OpenID Connect.",
			},

			"use_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_OIDC", false),
				Description: "Allow OpenID Connect to be used for authentication",
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

			// Managed Tracking GUID for User-agent
			"partner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
				DefaultFunc:  schema.EnvDefaultFunc("ARM_PARTNER_ID", ""),
				Description:  "A GUID/UUID that is registered with Microsoft to facilitate partner resource usage attribution.",
			},

			"disable_correlation_request_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_DISABLE_CORRELATION_REQUEST_ID", false),
				Description: "This will disable the x-ms-correlation-request-id header.",
			},

			"disable_terraform_partner_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_DISABLE_TERRAFORM_PARTNER_ID", false),
				Description: "This will disable the Terraform Partner ID which is used if a custom `partner_id` isn't specified.",
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
		// var auxTenants []string
		// if v, ok := d.Get("auxiliary_tenant_ids").([]interface{}); ok && len(v) > 0 {
		// 	auxTenants = *utils.ExpandStringSlice(v)
		// } else if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
		// 	auxTenants = strings.Split(v, ";")
		// }

		var cloudConfig cloud.Configuration
		env := d.Get("environment").(string)
		switch strings.ToLower(env) {
		case "public":
			cloudConfig = cloud.AzurePublic
		case "usgovernment":
			cloudConfig = cloud.AzureGovernment
		case "china":
			cloudConfig = cloud.AzureChina
		default:
			return nil, diag.Errorf("unknown `environment` specified: %q", env)
		}

		// Maps the auth related environment variables used in the provider to what azidentity honors.
		if v := d.Get("tenant_id").(string); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_TENANT_ID", v)
		}
		if v := d.Get("client_id").(string); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_CLIENT_ID", v)
		}
		if v := d.Get("client_secret").(string); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_CLIENT_SECRET", v)
		}
		if v := d.Get("client_certificate_path").(string); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", v)
		}
		if v := d.Get("client_certificate_password").(string); len(v) != 0 {
			// #nosec G104
			os.Setenv("AZURE_CLIENT_CERTIFICATE_PASSWORD", v)
		}

		getAssertion := func(ctx context.Context) (string, error) {
			if v := d.Get("oidc_token").(string); len(v) != 0 {
				return v, nil
			}
			if v := d.Get("oidc_request_url").(string); len(v) != 0 {
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, v, http.NoBody)
				if err != nil {
					return "", fmt.Errorf("githubAssertion: failed to build request")
				}
				req.Header.Set("Accept", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.Get("oidc_request_token").(string)))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return "", fmt.Errorf("githubAssertion: cannot request token: %v", err)
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
				if err != nil {
					return "", fmt.Errorf("githubAssertion: cannot parse response: %v", err)
				}

				if c := resp.StatusCode; c < 200 || c > 299 {
					return "", fmt.Errorf("githubAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
				}

				var tokenRes struct {
					Count *int    `json:"count"`
					Value *string `json:"value"`
				}
				if err := json.Unmarshal(body, &tokenRes); err != nil {
					return "", fmt.Errorf("githubAssertion: cannot unmarshal response: %v", err)
				}

				return *tokenRes.Value, nil
			}
			return "", fmt.Errorf("OIDC: failed getting token")
		}

		oidcCredentialConstructorErrorHandler := func(numberOfSuccessfulCredentials int, errorMessages []string) (err error) {
			errorMessage := strings.Join(errorMessages, "\n\t")

			if numberOfSuccessfulCredentials == 0 {
				return fmt.Errorf(errorMessage)
			}

			return nil
		}

		var creds []azcore.TokenCredential
		var errorMessages []string

		if d.Get("use_oidc").(bool) {
			if v := d.Get("oidc_token_file_path").(string); len(v) != 0 {
				// #nosec G104
				os.Setenv("AZURE_FEDERATED_TOKEN_FILE", v)
			}
			v := d.Get("oidc_authority_host").(string)
			// #nosec G104
			os.Setenv("AZURE_AUTHORITY_HOST", v)

			if v := d.Get("tenant_id").(string); len(v) != 0 {
				if v := d.Get("client_id").(string); len(v) != 0 {
					w := d.Get("oidc_request_url").(string)
					if v := d.Get("oidc_token").(string); len(v) != 0 || len(w) != 0 {
						oidcCred, err := azidentity.NewClientAssertionCredential(d.Get("tenant_id").(string), d.Get("client_id").(string), getAssertion, &azidentity.ClientAssertionCredentialOptions{
							ClientOptions: azcore.ClientOptions{},
						})
						if err == nil {
							creds = append(creds, oidcCred)
						} else {
							errorMessages = append(errorMessages, "oidcCredential: "+err.Error())
						}
					}
				}
			}
		}

		defaultCred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: cloudConfig,
			},
			TenantID: d.Get("tenant_id").(string),
		})
		if err == nil {
			creds = append(creds, defaultCred)
		} else {
			errorMessages = append(errorMessages, "defaultCredential: "+err.Error())
		}

		err = oidcCredentialConstructorErrorHandler(len(creds), errorMessages)
		if err != nil {
			return nil, diag.Errorf("failed to obtain a credential: %v", err)
		}

		chain, err := azidentity.NewChainedTokenCredential(creds, nil)
		if err != nil {
			return nil, diag.Errorf("%v", err)
		}

		cred := &chainTokenCredential{chain: chain}

		copt := &clients.Option{
			SubscriptionId:       d.Get("subscription_id").(string),
			Cred:                 cred,
			CloudCfg:             cloudConfig,
			ApplicationUserAgent: buildUserAgent(p.TerraformVersion, d.Get("partner_id").(string), d.Get("disable_terraform_partner_id").(bool)),
			Features: features.UserFeatures{
				DefaultTags:     tags.ExpandTags(d.Get("default_tags").(map[string]interface{})),
				DefaultLocation: location.Normalize(d.Get("default_location").(string)),
			},
			SkipProviderRegistration:    d.Get("skip_provider_registration").(bool),
			DisableCorrelationRequestID: d.Get("disable_correlation_request_id").(bool),
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

func buildUserAgent(terraformVersion string, partnerID string, disableTerraformPartnerID bool) string {
	if terraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		terraformVersion = "0.11+compatible"
	}

	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", terraformVersion, meta.SDKVersionString())
	providerUserAgent := fmt.Sprintf("terraform-provider-azapi/%s", version.ProviderVersion)
	userAgent := strings.TrimSpace(fmt.Sprintf("%s %s", tfUserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, azureAgent)
	}

	// only one pid can be interpreted currently
	// hence, send partner ID if present, otherwise send Terraform GUID
	// unless users have opted out
	if partnerID == "" && !disableTerraformPartnerID {
		// Microsoft’s Terraform Partner ID is this specific GUID
		partnerID = "222c6c49-1b0a-5959-a213-6608f9eb8820"
	}

	if partnerID != "" {
		userAgent = fmt.Sprintf("%s pid-%s", userAgent, partnerID)
	}
	return userAgent
}

type chainTokenCredential struct {
	chain *azidentity.ChainedTokenCredential
}

func (c chainTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return c.chain.GetToken(ctx, opts)
}
