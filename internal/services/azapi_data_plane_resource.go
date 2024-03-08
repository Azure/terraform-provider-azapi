package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/services/defaults"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type DataPlaneResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	ParentID              types.String   `tfsdk:"parent_id"`
	Type                  types.String   `tfsdk:"type"`
	Body                  types.String   `tfsdk:"body"`
	IgnoreCasing          types.Bool     `tfsdk:"ignore_casing"`
	IgnoreMissingProperty types.Bool     `tfsdk:"ignore_missing_property"`
	ResponseExportValues  types.List     `tfsdk:"response_export_values"`
	Locks                 types.List     `tfsdk:"locks"`
	Output                types.String   `tfsdk:"output"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

type DataPlaneResource struct {
	ProviderData *clients.Client
}

var _ resource.Resource = &DataPlaneResource{}
var _ resource.ResourceWithConfigure = &DataPlaneResource{}
var _ resource.ResourceWithModifyPlan = &DataPlaneResource{}

func (r *DataPlaneResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *DataPlaneResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_data_plane_resource"
}

func (r *DataPlaneResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"parent_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
			},

			"body": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.StringDefault("{}"),
				Validators: []validator.String{
					myvalidator.StringIsJSON(),
				},
				PlanModifiers: []planmodifier.String{
					myplanmodifier.UseStateWhen(func(a, b types.String) bool {
						return utils.NormalizeJson(a.ValueString()) == utils.NormalizeJson(b.ValueString())
					}),
				},
			},

			"ignore_casing": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.BoolDefault(false),
			},

			"ignore_missing_property": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.BoolDefault(true),
			},

			"response_export_values": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
			},

			"locks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
			},

			"output": schema.StringAttribute{
				Computed: true,
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
	}
}

func (r *DataPlaneResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	var config, plan, state *DataPlaneResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// destroy doesn't need to modify plan
	if config == nil {
		return
	}

	if state == nil || !plan.ResponseExportValues.Equal(state.ResponseExportValues) || utils.NormalizeJson(plan.Body.ValueString()) != utils.NormalizeJson(state.Body.ValueString()) {
		plan.Output = types.StringUnknown()
	}
}

func (r *DataPlaneResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	r.CreateUpdate(ctx, request.Plan, &response.State, &response.Diagnostics)
}

func (r *DataPlaneResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	r.CreateUpdate(ctx, request.Plan, &response.State, &response.Diagnostics)
}

func (r *DataPlaneResource) CreateUpdate(ctx context.Context, plan tfsdk.Plan, state *tfsdk.State, diagnostics *diag.Diagnostics) {
	var model DataPlaneResourceModel
	if diagnostics.Append(plan.Get(ctx, &model)...); diagnostics.HasError() {
		return
	}

	createUpdateTimeout, diags := model.Timeouts.Create(ctx, 30*time.Minute)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createUpdateTimeout)
	defer cancel()

	id, err := parse.NewDataPlaneResourceId(model.Name.ValueString(), model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	client := r.ProviderData.DataPlaneClient
	if isNewResource := state == nil || state.Raw.IsNull(); isNewResource {
		_, err = client.Get(ctx, id)
		if err == nil {
			diagnostics.AddError("Resource already exists", tf.ImportAsExistsError("azapi_data_plane_resource", id.ID()).Error())
			return
		}
		if !utils.ResponseErrorWasNotFound(err) {
			diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("checking for presence of existing %s: %+v", id, err).Error())
			return
		}
	}

	body := make(map[string]interface{})
	err = json.Unmarshal([]byte(model.Body.ValueString()), &body)
	if err != nil {
		diagnostics.AddError("Invalid configuration", fmt.Sprintf(`The argument "body" is invalid, value: %s, error: %s`, model.Body.ValueString(), err))
		return
	}

	for _, id := range AsStringList(model.Locks) {
		locks.ByID(id)
		defer locks.UnlockByID(id)
	}

	responseBody, err := client.CreateOrUpdateThenPoll(ctx, id, body)
	if err != nil {
		diagnostics.AddError("Failed to create/update resource", fmt.Errorf("creating/updating %q: %+v", id, err).Error())
		return
	}

	model.ID = basetypes.NewStringValue(id.ID())
	model.Output = basetypes.NewStringValue(flattenOutput(responseBody, AsStringList(model.ResponseExportValues)))

	diagnostics.Append(state.Set(ctx, model)...)
}

func (r *DataPlaneResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var model *DataPlaneResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := model.Timeouts.Read(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	id, err := parse.DataPlaneResourceIDWithResourceType(model.ID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	client := r.ProviderData.DataPlaneClient
	responseBody, err := client.Get(ctx, id)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("[INFO] Error reading %q - removing from state", id.ID()))
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
		return
	}

	bodyJson := model.Body.ValueString()
	requestBody := make(map[string]interface{})
	err = json.Unmarshal([]byte(bodyJson), &requestBody)
	if err != nil && bodyJson != "" {
		response.Diagnostics.AddError("Invalid configuration", fmt.Sprintf(`The argument "body" is invalid, value: %s, error: %s`, bodyJson, err))
		return
	}

	option := utils.UpdateJsonOption{
		IgnoreCasing:          model.IgnoreCasing.ValueBool(),
		IgnoreMissingProperty: model.IgnoreMissingProperty.ValueBool(),
	}
	data, err := json.Marshal(utils.UpdateObject(requestBody, responseBody, option))
	if err != nil {
		response.Diagnostics.AddError("Invalid body", err.Error())
		return
	}

	model.Body = basetypes.NewStringValue(string(data))
	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.Type = basetypes.NewStringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))
	model.Output = basetypes.NewStringValue(flattenOutput(responseBody, AsStringList(model.ResponseExportValues)))

	response.Diagnostics.Append(response.State.Set(ctx, model)...)
}

func (r *DataPlaneResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	client := r.ProviderData.DataPlaneClient

	var model *DataPlaneResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := model.Timeouts.Delete(ctx, 30*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	id, err := parse.DataPlaneResourceIDWithResourceType(model.ID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	for _, lockId := range AsStringList(model.Locks) {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	_, err = client.DeleteThenPoll(ctx, id)
	if err != nil && !utils.ResponseErrorWasNotFound(err) {
		response.Diagnostics.AddError("Failed to delete resource", fmt.Errorf("deleting %s: %+v", id, err).Error())
	}
}
