package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/customization"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/ephemeral/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DataPlaneResourceEphemeral struct {
	ProviderData *clients.Client
}

var (
	_ ephemeral.EphemeralResource                   = &DataPlaneResourceEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure      = &DataPlaneResourceEphemeral{}
	_ ephemeral.EphemeralResourceWithValidateConfig = &DataPlaneResourceEphemeral{}
)

func (r *DataPlaneResourceEphemeral) Configure(ctx context.Context, request ephemeral.ConfigureRequest, response *ephemeral.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring azapi_data_plane_resource ephemeral resource")
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

type DataPlaneResourceEphemeralModel struct {
	ID                   types.String     `tfsdk:"id"`
	Name                 types.String     `tfsdk:"name"`
	ParentID             types.String     `tfsdk:"parent_id"`
	Type                 types.String     `tfsdk:"type"`
	Body                 types.Dynamic    `tfsdk:"body"`
	ResponseExportValues types.Dynamic    `tfsdk:"response_export_values"`
	Output               types.Dynamic    `tfsdk:"output"`
	Timeouts             timeouts.Value   `tfsdk:"timeouts"`
	Retry                retry.RetryValue `tfsdk:"retry"`
}

func (r *DataPlaneResourceEphemeral) Metadata(ctx context.Context, request ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_data_plane_resource"
}

func (r *DataPlaneResourceEphemeral) Schema(ctx context.Context, request ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This ephemeral resource can read the properties of an Azure data plane resource without persisting them to Terraform state. Ephemeral resources are re-read on every operation, making them suitable for retrieving sensitive or short-lived data plane values.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.ID(),
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specifies the name (identifier segment) of the data plane resource.",
			},
			"parent_id": schema.StringAttribute{
				Required:            true,
				Validators:          []validator.String{myvalidator.StringIsNotEmpty()},
				MarkdownDescription: "The ID of the azure resource in which this resource exists.",
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.DataPlaneType(),
			},
			"body": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.BodyResponse(),
			},
			"response_export_values": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: docstrings.ResponseExportValues(),
			},
			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("ephemeral.azapi_data_plane_resource"),
			},
			"retry": retry.RetryEphemeralSchema(ctx),
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx),
		},
	}
}

func (r *DataPlaneResourceEphemeral) ValidateConfig(ctx context.Context, request ephemeral.ValidateConfigRequest, response *ephemeral.ValidateConfigResponse) {
	var config *DataPlaneResourceEphemeralModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}
	if config == nil {
		return
	}
	if err := validateDataPlaneResourceName(&DataPlaneResourceModel{
		Name:     config.Name,
		ParentID: config.ParentID,
		Type:     config.Type,
	}); err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
	}
}

func (r *DataPlaneResourceEphemeral) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {
	var model *DataPlaneResourceEphemeralModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	openTimeout, diags := model.Timeouts.Open(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, openTimeout)
	defer cancel()

	id, err := parse.NewDataPlaneResourceId(model.Name.ValueString(), model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}
	ctx = tflog.SetField(ctx, "resource_id", id.ID())

	client := r.ProviderData.DataPlaneClient
	requestOptions := clients.RequestOptions{}
	requestOptions.RetryOptions, requestOptions.LastRetryError = clients.NewRetryOptions(model.Retry)

	var responseBody interface{}
	if customizedResource := customization.GetCustomization(model.Type.ValueString()); customizedResource != nil && (*customizedResource).ReadFunc() != nil {
		responseBody, err = (*customizedResource).ReadFunc()(ctx, *r.ProviderData, id, requestOptions)
	} else {
		responseBody, err = client.Get(ctx, id, requestOptions)
	}

	if err != nil {
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
		return
	}

	data, err := json.Marshal(responseBody)
	if err != nil {
		response.Diagnostics.AddError("Invalid body", err.Error())
		return
	}

	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, nil)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	payload, err := dynamic.FromJSONImplied(data)
	if err != nil {
		response.Diagnostics.AddError("Invalid payload", err.Error())
		return
	}
	model.Body = payload

	model.ID = basetypes.NewStringValue(id.ID())
	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.Type = basetypes.NewStringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

	response.Diagnostics.Append(response.Result.Set(ctx, model)...)
}
