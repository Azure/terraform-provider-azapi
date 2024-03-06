package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/version"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func AzureProvider() provider.Provider {
	return &Provider{}
}

type Provider struct {
}

type providerData struct {
	SubscriptionID              types.String `tfsdk:"subscription_id"`
	ClientID                    types.String `tfsdk:"client_id"`
	TenantID                    types.String `tfsdk:"tenant_id"`
	AuxiliaryTenantIDs          types.List   `tfsdk:"auxiliary_tenant_ids"`
	Endpoint                    types.List   `tfsdk:"endpoint"`
	Environment                 types.String `tfsdk:"environment"`
	ClientCertificatePath       types.String `tfsdk:"client_certificate_path"`
	ClientCertificatePassword   types.String `tfsdk:"client_certificate_password"`
	ClientSecret                types.String `tfsdk:"client_secret"`
	SkipProviderRegistration    types.Bool   `tfsdk:"skip_provider_registration"`
	OIDCRequestToken            types.String `tfsdk:"oidc_request_token"`
	OIDCRequestURL              types.String `tfsdk:"oidc_request_url"`
	OIDCToken                   types.String `tfsdk:"oidc_token"`
	OIDCTokenFilePath           types.String `tfsdk:"oidc_token_file_path"`
	UseOIDC                     types.Bool   `tfsdk:"use_oidc"`
	UseCLI                      types.Bool   `tfsdk:"use_cli"`
	UseMSI                      types.Bool   `tfsdk:"use_msi"`
	PartnerID                   types.String `tfsdk:"partner_id"`
	CustomCorrelationRequestID  types.String `tfsdk:"custom_correlation_request_id"`
	DisableCorrelationRequestID types.Bool   `tfsdk:"disable_correlation_request_id"`
	DisableTerraformPartnerID   types.Bool   `tfsdk:"disable_terraform_partner_id"`
	DefaultName                 types.String `tfsdk:"default_name"`
	DefaultNamingPrefix         types.String `tfsdk:"default_naming_prefix"`
	DefaultNamingSuffix         types.String `tfsdk:"default_naming_suffix"`
	DefaultLocation             types.String `tfsdk:"default_location"`
	DefaultTags                 types.Map    `tfsdk:"default_tags"`
}

type providerEndpointData struct {
	ActiveDirectoryAuthorityHost types.String `tfsdk:"active_directory_authority_host"`
	ResourceManagerEndpoint      types.String `tfsdk:"resource_manager_endpoint"`
	ResourceManagerAudience      types.String `tfsdk:"resource_manager_audience"`
}

func (p Provider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "azapi"
}

func (p Provider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "The Azure API Provider",
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Subscription ID which should be used.",
			},

			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Client ID which should be used.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Tenant ID which should be used.",
			},

			"auxiliary_tenant_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators:  []validator.List{listvalidator.SizeAtMost(3)},
				Description: "The Auxiliary Tenant IDs which should be used.",
			},

			"endpoint": schema.ListNestedAttribute{
				Optional:   true,
				Validators: []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active_directory_authority_host": schema.StringAttribute{
							Optional:    true,
							Description: "The Active Directory login endpoint which should be used.",
						},

						"resource_manager_endpoint": schema.StringAttribute{
							Optional:    true,
							Description: "The Resource Manager Endpoint which should be used.",
						},

						"resource_manager_audience": schema.StringAttribute{
							Optional:    true,
							Description: "The resource ID to obtain AD tokens for.",
						},
					},
				},
			},

			"environment": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("public", "usgovernment", "china"),
				},
				Description: "The Cloud Environment which should be used. Possible values are public, usgovernment and china. Defaults to public.",
			},

			// TODO@mgd: the metadata_host is used to retrieve metadata from Azure to identify current environment, this is used to eliminate Azure Stack usage, in which case the provider doesn't support.
			// "metadata_host": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	DefaultFunc: schema.EnvDefaultFunc("ARM_METADATA_HOSTNAME", ""),
			// 	Description: "The Hostname which should be used for the Azure Metadata Service.",
			// },

			// Client Certificate specific fields
			"client_certificate_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			"client_certificate_password": schema.StringAttribute{
				Optional:    true,
				Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"skip_provider_registration": schema.BoolAttribute{
				Optional:    true,
				Description: "Should the Provider skip registering all of the Resource Providers that it supports, if they're not already registered?",
			},

			// OIDC specific fields
			"oidc_request_token": schema.StringAttribute{
				Optional:    true,
				Description: "The bearer token for the request to the OIDC provider. For use When authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_request_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL for the OIDC provider from which to request an ID token. For use When authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token": schema.StringAttribute{
				Optional:    true,
				Description: "The OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token_file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a file containing an OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"use_oidc": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow OpenID Connect to be used for authentication",
			},

			// Azure CLI specific fields
			"use_cli": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow Azure CLI to be used for Authentication.",
			},

			// Managed Service Identity specific fields
			"use_msi": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow Managed Service Identity to be used for Authentication.",
			},

			// TODO@mgd: azidentity doesn't support msi_endpoint
			// "msi_endpoint": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
			// 	Description: "The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. ",
			// },

			// Managed Tracking GUID for User-agent
			"partner_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.Any(myvalidator.StringIsUUID(), myvalidator.StringIsEmpty()),
				},
				Description: "A GUID/UUID that is registered with Microsoft to facilitate partner resource usage attribution.",
			},

			"custom_correlation_request_id": schema.StringAttribute{
				Optional:    true,
				Description: "The value of the x-ms-correlation-request-id header (otherwise an auto-generated UUID will be used).",
			},

			"disable_correlation_request_id": schema.BoolAttribute{
				Optional:    true,
				Description: "This will disable the x-ms-correlation-request-id header.",
			},

			"disable_terraform_partner_id": schema.BoolAttribute{
				Optional:    true,
				Description: "This will disable the Terraform Partner ID which is used if a custom `partner_id` isn't specified.",
			},

			"default_name": schema.StringAttribute{
				Optional:    true,
				Description: "The default name which should be used for resources.",
			},

			"default_naming_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "The default prefix which should be used for resources.",
			},

			"default_naming_suffix": schema.StringAttribute{
				Optional:    true,
				Description: "The default suffix which should be used for resources.",
			},

			"default_location": schema.StringAttribute{
				Optional:    true,
				Description: "The default location which should be used for resources.",
			},

			"default_tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{
					tags.Validator(),
				},
				Description: "The default tags which should be used for resources.",
			},
		},
	}
}

func (p Provider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var model providerData
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	if !model.DefaultName.IsNull() && (!model.DefaultNamingPrefix.IsNull() || !model.DefaultNamingSuffix.IsNull()) {
		response.Diagnostics.AddError("Invalid `default_name` value.", "The `default_name` value cannot be used with `default_naming_prefix` or `default_naming_suffix`.")
		return
	}

	// set the defaults from environment variables
	if model.SubscriptionID.IsNull() {
		if v := os.Getenv("ARM_SUBSCRIPTION_ID"); v != "" {
			model.SubscriptionID = types.StringValue(v)
		}
	}
	if model.ClientID.IsNull() {
		if v := os.Getenv("ARM_CLIENT_ID"); v != "" {
			model.ClientID = types.StringValue(v)
		}
	}

	if model.TenantID.IsNull() {
		if v := os.Getenv("ARM_TENANT_ID"); v != "" {
			model.TenantID = types.StringValue(v)
		}
	}

	if model.Endpoint.IsNull() {
		activeDirectoryAuthorityHost := os.Getenv("ARM_ACTIVE_DIRECTORY_AUTHORITY_HOST")
		resourceManagerEndpoint := os.Getenv("ARM_RESOURCE_MANAGER_ENDPOINT")
		resourceManagerAudience := os.Getenv("ARM_RESOURCE_MANAGER_AUDIENCE")
		attrTypes := make(map[string]attr.Type)
		attrTypes["active_directory_authority_host"] = types.StringType
		attrTypes["resource_manager_endpoint"] = types.StringType
		attrTypes["resource_manager_audience"] = types.StringType
		model.Endpoint = types.ListValueMust(types.ObjectType{
			AttrTypes: attrTypes,
		}, []attr.Value{
			types.ObjectValueMust(attrTypes, map[string]attr.Value{
				"active_directory_authority_host": types.StringValue(activeDirectoryAuthorityHost),
				"resource_manager_endpoint":       types.StringValue(resourceManagerEndpoint),
				"resource_manager_audience":       types.StringValue(resourceManagerAudience),
			}),
		})
	}

	if model.Environment.IsNull() {
		if v := os.Getenv("ARM_ENVIRONMENT"); v != "" {
			model.Environment = types.StringValue(v)
		} else {
			model.Environment = types.StringValue("public")
		}
	}

	if model.AuxiliaryTenantIDs.IsNull() {
		if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
			values := make([]attr.Value, 0)
			for _, v := range strings.Split(v, ";") {
				values = append(values, types.StringValue(v))
			}
			model.AuxiliaryTenantIDs = types.ListValueMust(types.StringType, values)
		}
	}

	if model.ClientCertificatePath.IsNull() {
		if v := os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"); v != "" {
			model.ClientCertificatePath = types.StringValue(v)
		}
	}

	if model.ClientCertificatePassword.IsNull() {
		if v := os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"); v != "" {
			model.ClientCertificatePassword = types.StringValue(v)
		}
	}

	if model.ClientSecret.IsNull() {
		if v := os.Getenv("ARM_CLIENT_SECRET"); v != "" {
			model.ClientSecret = types.StringValue(v)
		}
	}

	if model.SkipProviderRegistration.IsNull() {
		if v := os.Getenv("ARM_SKIP_PROVIDER_REGISTRATION"); v != "" {
			model.SkipProviderRegistration = types.BoolValue(v == "true")
		} else {
			model.SkipProviderRegistration = types.BoolValue(false)
		}
	}

	if model.OIDCRequestToken.IsNull() {
		if v := os.Getenv("ARM_OIDC_REQUEST_TOKEN"); v != "" {
			model.OIDCRequestToken = types.StringValue(v)
		} else if v := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"); v != "" {
			model.OIDCRequestToken = types.StringValue(v)
		}
	}

	if model.OIDCRequestURL.IsNull() {
		if v := os.Getenv("ARM_OIDC_REQUEST_URL"); v != "" {
			model.OIDCRequestURL = types.StringValue(v)
		} else if v := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"); v != "" {
			model.OIDCRequestURL = types.StringValue(v)
		}
	}

	if model.OIDCToken.IsNull() {
		if v := os.Getenv("ARM_OIDC_TOKEN"); v != "" {
			model.OIDCToken = types.StringValue(v)
		}
	}

	if model.OIDCTokenFilePath.IsNull() {
		if v := os.Getenv("ARM_OIDC_TOKEN_FILE_PATH"); v != "" {
			model.OIDCTokenFilePath = types.StringValue(v)
		}
	}

	if model.UseOIDC.IsNull() {
		if v := os.Getenv("ARM_USE_OIDC"); v != "" {
			model.UseOIDC = types.BoolValue(v == "true")
		} else {
			model.UseOIDC = types.BoolValue(false)
		}
	}

	if model.UseCLI.IsNull() {
		if v := os.Getenv("ARM_USE_CLI"); v != "" {
			model.UseCLI = types.BoolValue(v == "true")
		} else {
			model.UseCLI = types.BoolValue(true)
		}
	}

	if model.UseMSI.IsNull() {
		if v := os.Getenv("ARM_USE_MSI"); v != "" {
			model.UseMSI = types.BoolValue(v == "true")
		} else {
			model.UseMSI = types.BoolValue(true)
		}
	}

	if model.PartnerID.IsNull() {
		if v := os.Getenv("ARM_PARTNER_ID"); v != "" {
			model.PartnerID = types.StringValue(v)
		}
	}

	if model.CustomCorrelationRequestID.IsNull() {
		if v := os.Getenv("ARM_CORRELATION_REQUEST_ID"); v != "" {
			model.CustomCorrelationRequestID = types.StringValue(v)
		}
	}

	if model.DisableCorrelationRequestID.IsNull() {
		if v := os.Getenv("ARM_DISABLE_CORRELATION_REQUEST_ID"); v != "" {
			model.DisableCorrelationRequestID = types.BoolValue(v == "true")
		} else {
			model.DisableCorrelationRequestID = types.BoolValue(false)
		}
	}

	if model.DisableTerraformPartnerID.IsNull() {
		if v := os.Getenv("ARM_DISABLE_TERRAFORM_PARTNER_ID"); v != "" {
			model.DisableTerraformPartnerID = types.BoolValue(v == "true")
		} else {
			model.DisableTerraformPartnerID = types.BoolValue(false)
		}
	}

	var cloudConfig cloud.Configuration
	env := model.Environment.ValueString()
	switch strings.ToLower(env) {
	case "public":
		cloudConfig = cloud.AzurePublic
	case "usgovernment":
		cloudConfig = cloud.AzureGovernment
	case "china":
		cloudConfig = cloud.AzureChina
	default:
		response.Diagnostics.AddError("Invalid `environment` value.", fmt.Sprintf("The `environment` value '%s' is invalid. Valid values are 'public', 'usgovernment' and 'china'.", env))
		return
	}

	if elements := model.Endpoint.Elements(); len(elements) != 0 {
		var endpoint providerEndpointData
		diags := elements[0].(basetypes.ObjectValue).As(ctx, &endpoint, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		response.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		resourceManagerEndpoint := cloudConfig.Services[cloud.ResourceManager].Endpoint
		resourceManagerAudience := cloudConfig.Services[cloud.ResourceManager].Audience
		if v := endpoint.ResourceManagerEndpoint.ValueString(); v != "" {
			resourceManagerEndpoint = v
		}
		if v := endpoint.ResourceManagerAudience.ValueString(); v != "" {
			resourceManagerAudience = v
		}
		cloudConfig.Services[cloud.ResourceManager] = cloud.ServiceConfiguration{
			Endpoint: resourceManagerEndpoint,
			Audience: resourceManagerAudience,
		}
		if v := endpoint.ActiveDirectoryAuthorityHost.ValueString(); v != "" {
			cloudConfig.ActiveDirectoryAuthorityHost = v
		}
	}

	// Maps the auth related environment variables used in the provider to what azidentity honors.
	if v := model.TenantID.ValueString(); v != "" {
		_ = os.Setenv("AZURE_TENANT_ID", v)
	}
	if v := model.ClientID.ValueString(); v != "" {
		_ = os.Setenv("AZURE_CLIENT_ID", v)
	}
	if v := model.ClientSecret.ValueString(); v != "" {
		_ = os.Setenv("AZURE_CLIENT_SECRET", v)
	}
	if v := model.ClientCertificatePath.ValueString(); v != "" {
		_ = os.Setenv("AZURE_CERTIFICATE_PATH", v)
	}
	if v := model.ClientCertificatePassword.ValueString(); v != "" {
		_ = os.Setenv("AZURE_CERTIFICATE_PASSWORD", v)
	}
	var auxTenants []string
	if elements := model.AuxiliaryTenantIDs.Elements(); len(elements) != 0 {
		for _, element := range elements {
			auxTenants = append(auxTenants, element.(basetypes.StringValue).ValueString())
		}
		_ = os.Setenv("AZURE_AUXILIARY_TENANT_IDS", strings.Join(auxTenants, ";"))
	}

	option := &azidentity.DefaultAzureCredentialOptions{
		AdditionallyAllowedTenants: auxTenants,
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
		TenantID: model.TenantID.ValueString(),
	}

	cred, err := newDefaultAzureCredential(model, option)
	if err != nil {
		response.Diagnostics.AddError("Failed to obtain a credential.", err.Error())
		return
	}

	copt := &clients.Option{
		Cred:                 cred,
		CloudCfg:             cloudConfig,
		ApplicationUserAgent: buildUserAgent(request.TerraformVersion, model.PartnerID.ValueString(), model.DisableTerraformPartnerID.ValueBool()),
		Features: features.UserFeatures{
			DefaultTags:         tags.ExpandTags(model.DefaultTags),
			DefaultLocation:     location.Normalize(model.DefaultLocation.ValueString()),
			DefaultNaming:       model.DefaultName.ValueString(),
			DefaultNamingPrefix: model.DefaultNamingPrefix.ValueString(),
			DefaultNamingSuffix: model.DefaultNamingSuffix.ValueString(),
		},
		SkipProviderRegistration:    model.SkipProviderRegistration.ValueBool(),
		DisableCorrelationRequestID: model.DisableCorrelationRequestID.ValueBool(),
		CustomCorrelationRequestID:  model.CustomCorrelationRequestID.ValueString(),
		SubscriptionId:              model.SubscriptionID.ValueString(),
	}

	client := &clients.Client{}
	if err = client.Build(ctx, copt); err != nil {
		response.Diagnostics.AddError("Error Building Client", err.Error())
		return
	}

	// load schema
	azure.GetAzureSchema()

	response.ResourceData = client
	response.DataSourceData = client
}

func (p Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &services.ResourceIdDataSource{}
		},
		func() datasource.DataSource {
			return &services.ResourceListDataSource{}
		},
		func() datasource.DataSource {
			return &services.ResourceActionDataSource{}
		},
		func() datasource.DataSource {
			return &services.AzapiResourceDataSource{}
		},
	}

}

func (p Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &services.AzapiResource{}
		},
		func() resource.Resource {
			return &services.AzapiUpdateResource{}
		},
		func() resource.Resource {
			return &services.ActionResource{}
		},
		func() resource.Resource {
			return &services.DataPlaneResource{}
		},
	}
}

func buildUserAgent(terraformVersion string, partnerID string, disableTerraformPartnerID bool) string {
	if terraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		terraformVersion = "0.11+compatible"
	}
	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io)", terraformVersion)
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

func newDefaultAzureCredential(model providerData, options *azidentity.DefaultAzureCredentialOptions) (*azidentity.ChainedTokenCredential, error) {
	var creds []azcore.TokenCredential

	if options == nil {
		options = &azidentity.DefaultAzureCredentialOptions{}
	}

	if model.UseOIDC.ValueBool() {
		oidcCred, err := NewOidcCredential(&OidcCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: options.Cloud,
			},
			AdditionallyAllowedTenants: options.AdditionallyAllowedTenants,
			TenantID:                   model.TenantID.ValueString(),
			ClientID:                   model.ClientID.ValueString(),
			RequestToken:               model.OIDCRequestToken.ValueString(),
			RequestUrl:                 model.OIDCRequestURL.ValueString(),
			Token:                      model.OIDCToken.ValueString(),
			TokenFilePath:              model.OIDCTokenFilePath.ValueString(),
		})

		if err == nil {
			creds = append(creds, oidcCred)
		} else {
			log.Printf("newDefaultAzureCredential failed to initialize oidc credential:\n\t%s", err.Error())
		}
	}

	envCred, err := azidentity.NewEnvironmentCredential(&azidentity.EnvironmentCredentialOptions{
		ClientOptions:            options.ClientOptions,
		DisableInstanceDiscovery: options.DisableInstanceDiscovery,
	})
	if err == nil {
		creds = append(creds, envCred)
	} else {
		log.Printf("newDefaultAzureCredential failed to initialize environment credential:\n\t%s", err.Error())
	}

	if model.UseMSI.ValueBool() {
		o := &azidentity.ManagedIdentityCredentialOptions{ClientOptions: options.ClientOptions}
		if ID, ok := os.LookupEnv("AZURE_CLIENT_ID"); ok {
			o.ID = azidentity.ClientID(ID)
		}
		miCred, err := NewManagedIdentityCredential(o)
		if err == nil {
			creds = append(creds, miCred)
		} else {
			log.Printf("newDefaultAzureCredential failed to initialize msi credential:\n\t%s", err.Error())
		}
	}

	if model.UseCLI.ValueBool() {
		cliCred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{
			AdditionallyAllowedTenants: options.AdditionallyAllowedTenants,
			TenantID:                   options.TenantID})
		if err == nil {
			creds = append(creds, cliCred)
		} else {
			log.Printf("newDefaultAzureCredential failed to initialize cli credential:\n\t%s", err.Error())
		}
	}

	if len(creds) == 0 {
		return nil, fmt.Errorf("no credentials were successfully initialized")
	}

	chain, err := azidentity.NewChainedTokenCredential(creds, nil)
	if err != nil {
		return nil, err
	}

	return chain, nil
}
