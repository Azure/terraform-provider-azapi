package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceListDataSourceModel struct {
	ID                   types.String        `tfsdk:"id"`
	Type                 types.String        `tfsdk:"type"`
	ParentID             types.String        `tfsdk:"parent_id"`
	ResponseExportValues types.Dynamic       `tfsdk:"response_export_values"`
	Output               types.Dynamic       `tfsdk:"output"`
	Timeouts             timeouts.Value      `tfsdk:"timeouts"`
	Retry                retry.RetryValue    `tfsdk:"retry"`
	Headers              map[string]string   `tfsdk:"headers"`
	QueryParameters      map[string][]string `tfsdk:"query_parameters"`
}

type ResourceListDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &ResourceListDataSource{}
var _ datasource.DataSourceWithConfigure = &ResourceListDataSource{}

func (r *ResourceListDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *ResourceListDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_list"
}

func (r *ResourceListDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
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

			"parent_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: docstrings.ParentID(),
			},

			"response_export_values": CommonAttributeResponseExportValues(),

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("data.azapi_resource_list"),
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

func (r *ResourceListDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model ResourceListDataSourceModel
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

	id, err := parse.NewResourceIDSkipScopeValidation("", model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	listUrl := strings.TrimSuffix(id.AzureResourceId, "/")

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

	responseBody, err := client.List(ctx, listUrl, id.ApiVersion, clients.NewRequestOptions(model.Headers, model.QueryParameters))
	if err != nil {
		response.Diagnostics.AddError("Failed to list resources", fmt.Sprintf("Failed to list resources, url: %s, error: %s", listUrl, err.Error()))
		return
	}

	model.ID = basetypes.NewStringValue(listUrl)
	var defaultOutput interface{}
	if !r.ProviderData.Features.DisableDefaultOutput {
		defaultOutput = responseBody
	}
	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, defaultOutput)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
