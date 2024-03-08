package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceListDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Type                 types.String `tfsdk:"type"`
	ParentID             types.String `tfsdk:"parent_id"`
	ResponseExportValues types.List   `tfsdk:"response_export_values"`
	Output               types.String `tfsdk:"output"`
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
				Computed: true,
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
			},

			"parent_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"response_export_values": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
			},

			"output": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ResourceListDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model ResourceListDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	id, err := parse.NewResourceIDSkipScopeValidation("", model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	listUrl := strings.TrimSuffix(id.AzureResourceId, "/")

	client := r.ProviderData.ResourceClient
	responseBody, err := client.List(ctx, listUrl, id.ApiVersion)
	if err != nil {
		response.Diagnostics.AddError("Failed to list resources", fmt.Sprintf("Failed to list resources, url: %s, error: %s", listUrl, err.Error()))
		return
	}

	model.ID = basetypes.NewStringValue(listUrl)
	model.Output = basetypes.NewStringValue(flattenOutput(responseBody, AsStringList(model.ResponseExportValues)))

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
