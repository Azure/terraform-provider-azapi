package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DataPlaneDataSourceModel struct {
	ID                            types.String     `tfsdk:"id"`
	Name                          types.String     `tfsdk:"name"`
	ParentID                      types.String     `tfsdk:"parent_id"`
	Type                          types.String     `tfsdk:"type"`
	ResponseExportValues          types.Dynamic    `tfsdk:"response_export_values"`
	SensitiveResponseExportValues types.Dynamic    `tfsdk:"sensitive_response_export_values"`
	Retry                         retry.RetryValue `tfsdk:"retry"`
	Locks                         types.List       `tfsdk:"locks"`
	Output                        types.Dynamic    `tfsdk:"output"`
	SensitiveOutput               types.Dynamic    `tfsdk:"sensitive_output"`
	Timeouts                      timeouts.Value   `tfsdk:"timeouts"`
	Headers                       types.Map        `tfsdk:"headers"`
	QueryParameters               types.Map        `tfsdk:"query_parameters"`
}

type DataPlaneDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &DataPlaneDataSource{}
var _ datasource.DataSourceWithConfigure = &DataPlaneDataSource{}

func (r *DataPlaneDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring azapi_data_plane_data_source")
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *DataPlaneDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_data_plane_resource"
}

func (r *DataPlaneDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This data source can read some Azure data plane resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.ID(),
			},

			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Specifies the name of the Azure data plane resource.",
			},

			"parent_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
				MarkdownDescription: "The ID of the Azure resource to get the data from.",
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type(),
			},

			"response_export_values": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: docstrings.ResponseExportValues(),
			},

			"sensitive_response_export_values": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: docstrings.SensitiveResponseExportValues(),
			},

			"retry": retry.RetrySchema(ctx),

			"locks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
				MarkdownDescription: docstrings.Locks(),
			},

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("data.azapi_data_plane_data_source"),
			},

			"sensitive_output": schema.DynamicAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: docstrings.SensitiveOutput("data.azapi_data_plane_data_source"),
			},

			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A mapping of headers to be sent with the read request.",
			},

			"query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the read request.",
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Read: true,
			}),
		},
	}
}

func (r *DataPlaneDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model *DataPlaneDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := model.Timeouts.Read(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	id, err := parse.NewDataPlaneResourceId(model.Name.ValueString(), model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}
	ctx = tflog.SetField(ctx, "resource_id", id.ID())

	// Ensure that the context deadline has been set before calling ConfigureClientWithCustomRetry().
	client := r.ProviderData.DataPlaneClient.ConfigureClientWithCustomRetry(ctx, model.Retry, false)

	responseBody, err := client.Get(ctx, id, clients.NewRequestOptions(AsMapOfString(model.Headers), AsMapOfLists(model.QueryParameters)))
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			response.Diagnostics.AddError("Resource not found", fmt.Errorf("resource %q not found", id).Error())
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("retrieving resource %q: %+v", id, err).Error())
		return
	}

	if _, err := json.Marshal(responseBody); err != nil {
		response.Diagnostics.AddError("Invalid body", err.Error())
		return
	}

	model.ID = basetypes.NewStringValue(id.ID())

	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, nil)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	sensitiveOutput, err := buildOutputFromBody(responseBody, model.SensitiveResponseExportValues, nil)
	if err != nil {
		response.Diagnostics.AddError("Failed to build sensitive output", err.Error())
		return
	}
	model.SensitiveOutput = sensitiveOutput

	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.Type = basetypes.NewStringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

	response.Diagnostics.Append(response.State.Set(ctx, model)...)
}
