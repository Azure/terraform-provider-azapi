package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceActionDataSourceModel struct {
	ID                   types.String        `tfsdk:"id"`
	ResourceID           types.String        `tfsdk:"resource_id"`
	Type                 types.String        `tfsdk:"type"`
	Action               types.String        `tfsdk:"action"`
	Method               types.String        `tfsdk:"method"`
	Body                 types.Dynamic       `tfsdk:"body"`
	ResponseExportValues types.Dynamic       `tfsdk:"response_export_values"`
	Output               types.Dynamic       `tfsdk:"output"`
	Timeouts             timeouts.Value      `tfsdk:"timeouts"`
	Retry                retry.RetryValue    `tfsdk:"retry"`
	Headers              map[string]string   `tfsdk:"headers"`
	QueryParameters      map[string][]string `tfsdk:"query_parameters"`
}

type ResourceActionDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &ResourceActionDataSource{}
var _ datasource.DataSourceWithConfigure = &ResourceActionDataSource{}

func (r *ResourceActionDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *ResourceActionDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_action"
}

func (r *ResourceActionDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
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

			"resource_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: "The ID of the Azure resource to perform the action on.",
			},

			"action": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: docstrings.ResourceAction(),
			},

			"method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("POST", "GET"),
				},
				MarkdownDescription: "The HTTP method to use when performing the action. Must be one of `POST`, `GET`. Defaults to `POST`.",
			},

			// The body attribute is a dynamic attribute that only allows users to specify the resource body as an HCL object
			"body": schema.DynamicAttribute{
				Optional: true,
				Validators: []validator.Dynamic{
					myvalidator.DynamicIsNotStringValidator(),
				},
			},

			"response_export_values": CommonAttributeResponseExportValues(),

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("data.azapi_resource_action"),
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

func (r *ResourceActionDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model ResourceActionDataSourceModel
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

	id, err := parse.ResourceIDWithResourceType(model.ResourceID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	var requestBody interface{}
	if err := unmarshalBody(model.Body, &requestBody); err != nil {
		response.Diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
		return
	}

	method := model.Method.ValueString()
	if method == "" {
		method = "POST"
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

	responseBody, err := client.Action(ctx, id.AzureResourceId, model.Action.ValueString(), id.ApiVersion, method, requestBody, clients.NewRequestOptions(model.Headers, model.QueryParameters))
	if err != nil {
		response.Diagnostics.AddError("Failed to perform action", fmt.Errorf("performing action %s of %q: %+v", model.Action.ValueString(), id, err).Error())
		return
	}

	model.ID = basetypes.NewStringValue(id.ID())

	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, nil)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
