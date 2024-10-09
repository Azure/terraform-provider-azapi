package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type AzapiResourceDataSourceModel struct {
	ID                   types.String        `tfsdk:"id"`
	Name                 types.String        `tfsdk:"name"`
	ParentID             types.String        `tfsdk:"parent_id"`
	ResourceID           types.String        `tfsdk:"resource_id"`
	Type                 types.String        `tfsdk:"type"`
	ResponseExportValues types.Dynamic       `tfsdk:"response_export_values"`
	Location             types.String        `tfsdk:"location"`
	Identity             types.List          `tfsdk:"identity"`
	Output               types.Dynamic       `tfsdk:"output"`
	Tags                 types.Map           `tfsdk:"tags"`
	Timeouts             timeouts.Value      `tfsdk:"timeouts"`
	Retry                retry.RetryValue    `tfsdk:"retry"`
	Headers              map[string]string   `tfsdk:"headers"`
	QueryParameters      map[string][]string `tfsdk:"query_parameters"`
}

type AzapiResourceDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &AzapiResourceDataSource{}
var _ datasource.DataSourceWithConfigure = &AzapiResourceDataSource{}

func (r *AzapiResourceDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *AzapiResourceDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource"
}

func (r *AzapiResourceDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This resource can access any existing Azure resource manager resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.ID(),
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type(),
			},

			`name`: schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
				MarkdownDescription: "Specifies the name of the Azure resource.",
			},

			"parent_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: docstrings.ParentID(),
			},

			"resource_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: "The ID of the Azure resource to retrieve.",
			},

			"location": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Location(),
			},

			"identity": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: docstrings.IdentityType(),
						},

						"principal_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: docstrings.IdentityPrincipalID(),
						},

						"tenant_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: docstrings.IdentityTenantID(),
						},

						"identity_ids": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: docstrings.IdentityIds(),
						},
					},
				},
			},

			"response_export_values": CommonAttributeResponseExportValues(),

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("data.azapi_resource"),
			},

			"tags": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "A mapping of tags which are assigned to the Azure resource.",
			},

			"retry": retry.SingleNestedAttribute(ctx),

			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of headers to include in the request",
			},

			"query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A map of query parameters to include in the request",
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Read: true,
			}),
		},
	}
}

func (r *AzapiResourceDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model AzapiResourceDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	model.Retry = model.Retry.AddDefaultValuesIfUnknownOrNull()

	readTimeout, diags := model.Timeouts.Read(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	var id parse.ResourceId
	resourceType := model.Type.ValueString()
	azureResourceType, _, _ := utils.GetAzureResourceTypeApiVersion(resourceType)
	if name := model.Name.ValueString(); len(name) != 0 {
		parentId := model.ParentID.ValueString()
		if parentId == "" && strings.EqualFold(azureResourceType, arm.ResourceGroupResourceType.String()) {
			parentId = fmt.Sprintf("/subscriptions/%s", r.ProviderData.Account.GetSubscriptionId())
		}
		buildId, err := parse.NewResourceID(name, parentId, resourceType)
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		id = buildId
	} else {
		resourceId := model.ResourceID.ValueString()
		if resourceId == "" && strings.EqualFold(azureResourceType, arm.SubscriptionResourceType.String()) {
			resourceId = fmt.Sprintf("/subscriptions/%s", r.ProviderData.Account.GetSubscriptionId())
		}
		buildId, err := parse.ResourceIDWithResourceType(resourceId, resourceType)
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		id = buildId
	}

	var client clients.Requester
	client = r.ProviderData.ResourceClient
	if !model.Retry.IsNull() && !model.Retry.IsUnknown() {
		bkof, regexps := clients.NewRetryableErrors(
			model.Retry.GetIntervalSeconds(),
			model.Retry.GetMaxIntervalSeconds(),
			model.Retry.GetMultiplier(),
			model.Retry.GetRandomizationFactor(),
			model.Retry.GetErrorMessageRegex(),
		)
		client = r.ProviderData.ResourceClient.WithRetry(bkof, regexps, nil, nil)
	}
	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(model.Headers, model.QueryParameters))
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			response.Diagnostics.AddError("Resource not found", fmt.Errorf("resource %q not found", id).Error())
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("retrieving resource %q: %+v", id, err).Error())
		return
	}

	model.ID = basetypes.NewStringValue(id.ID())
	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.ResourceID = basetypes.NewStringValue(id.AzureResourceId)
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		model.Tags = tags.FlattenTags(bodyMap["tags"])
		model.Location = basetypes.NewStringNull()
		if v := bodyMap["location"]; v != nil {
			model.Location = basetypes.NewStringValue(location.Normalize(v.(string)))
		}
		model.Identity = basetypes.NewListNull(identity.Model{}.ModelType())
		if v := identity.FlattenIdentity(bodyMap["identity"]); v != nil {
			model.Identity = identity.ToList(*v)
		}
	}

	var defaultOutput interface{}
	if !r.ProviderData.Features.DisableDefaultOutput {
		defaultOutput = id.ResourceDef.GetReadOnly(responseBody)
	}
	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, defaultOutput)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
