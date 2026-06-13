package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/Azure/terraform-provider-azapi/internal/services/customization"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DataPlaneResourceDataSourceModel struct {
	ID                   types.String     `tfsdk:"id"`
	Name                 types.String     `tfsdk:"name"`
	ParentID             types.String     `tfsdk:"parent_id"`
	Type                 types.String     `tfsdk:"type"`
	Identifiers          types.Map        `tfsdk:"identifiers"`
	Body                 types.Dynamic    `tfsdk:"body"`
	ResponseExportValues types.Dynamic    `tfsdk:"response_export_values"`
	Output               types.Dynamic    `tfsdk:"output"`
	Timeouts             timeouts.Value   `tfsdk:"timeouts"`
	Retry                retry.RetryValue `tfsdk:"retry"`
	Headers              types.Map        `tfsdk:"headers"`
	QueryParameters      types.Map        `tfsdk:"query_parameters"`
}

type DataPlaneResourceDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &DataPlaneResourceDataSource{}
var _ datasource.DataSourceWithConfigure = &DataPlaneResourceDataSource{}
var _ datasource.DataSourceWithValidateConfig = &DataPlaneResourceDataSource{}

func (r *DataPlaneResourceDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *DataPlaneResourceDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_data_plane_resource"
}

func (r *DataPlaneResourceDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This data source can read Azure data plane resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.ID(),
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specifies the name (identifier segment) of the data plane resource when the selected resource type uses a single `name` path segment.",
			},
			"parent_id": schema.StringAttribute{
				Required:            true,
				Validators:          []validator.String{myvalidator.StringIsNotEmpty()},
				MarkdownDescription: "The parent ID or endpoint prefix for the data plane resource being read.",
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type(),
			},
			"identifiers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A mapping of identifier placeholder values for data plane resource types that require multiple path identifiers, for example composite keys.",
			},
			"body": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Body(),
			},
			"response_export_values": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: docstrings.ResponseExportValues(),
			},
			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("data.azapi_data_plane_resource"),
			},
			"retry": retry.RetryDsSchema(ctx),
			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of headers to include in the request.",
			},
			"query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A map of query parameters to include in the request.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx),
		},
	}
}

func (r *DataPlaneResourceDataSource) ValidateConfig(ctx context.Context, request datasource.ValidateConfigRequest, response *datasource.ValidateConfigResponse) {
	var config *DataPlaneResourceDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}
	if config == nil {
		return
	}

	resourceConfig := &DataPlaneResourceModel{
		Name:        config.Name,
		ParentID:    config.ParentID,
		Type:        config.Type,
		Identifiers: config.Identifiers,
	}
	if err := validateDataPlaneResourceAddress(resourceConfig); err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
	}
}

func (r *DataPlaneResourceDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model DataPlaneResourceDataSourceModel
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

	id, err := parse.NewDataPlaneResourceIdWithIdentifiers(model.Name.ValueString(), model.ParentID.ValueString(), model.Type.ValueString(), common.AsMapOfString(model.Identifiers))
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}
	ctx = tflog.SetField(ctx, "resource_id", id.ID())

	requestOptions := clients.RequestOptions{
		Headers:         common.AsMapOfString(model.Headers),
		QueryParameters: clients.NewQueryParameters(common.AsMapOfLists(model.QueryParameters)),
		RetryOptions:    clients.NewRetryOptions(model.Retry),
	}

	var responseBody interface{}
	if customizedResource := customization.GetCustomization(model.Type.ValueString()); customizedResource != nil && (*customizedResource).ReadFunc() != nil {
		responseBody, err = (*customizedResource).ReadFunc()(ctx, *r.ProviderData, id, requestOptions)
	} else {
		responseBody, err = r.ProviderData.DataPlaneClient.Get(ctx, id, requestOptions)
	}
	if err != nil {
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
		return
	}

	bodyData, err := json.Marshal(responseBody)
	if err != nil {
		response.Diagnostics.AddError("Invalid body", err.Error())
		return
	}
	body, err := dynamic.FromJSONImplied(bodyData)
	if err != nil {
		response.Diagnostics.AddError("Invalid body", err.Error())
		return
	}

	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, responseBody)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}

	model.ID = basetypes.NewStringValue(id.ID())
	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.Type = basetypes.NewStringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))
	model.Identifiers = stringMapToTypesMap(id.Identifiers)
	model.Body = body
	model.Output = output

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
