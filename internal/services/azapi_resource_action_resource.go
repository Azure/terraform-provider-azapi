package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/services/defaults"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ActionResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Type                 types.String   `tfsdk:"type"`
	ResourceId           types.String   `tfsdk:"resource_id"`
	Action               types.String   `tfsdk:"action"`
	Method               types.String   `tfsdk:"method"`
	Body                 types.String   `tfsdk:"body"`
	Payload              types.Dynamic  `tfsdk:"payload"`
	When                 types.String   `tfsdk:"when"`
	Locks                types.List     `tfsdk:"locks"`
	ResponseExportValues types.List     `tfsdk:"response_export_values"`
	Output               types.String   `tfsdk:"output"`
	OutputPayload        types.Dynamic  `tfsdk:"output_payload"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type ActionResource struct {
	ProviderData *clients.Client
}

var _ resource.Resource = &ActionResource{}
var _ resource.ResourceWithConfigure = &ActionResource{}
var _ resource.ResourceWithModifyPlan = &ActionResource{}
var _ resource.ResourceWithValidateConfig = &ActionResource{}

func (r *ActionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *ActionResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_action"
}

func (r *ActionResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"resource_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"action": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.StringDefault("POST"),
				Validators: []validator.String{
					stringvalidator.OneOf("POST", "PATCH", "PUT", "DELETE", "GET", "HEAD"),
				},
			},

			"body": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					myvalidator.StringIsJSON(),
				},
				PlanModifiers: []planmodifier.String{
					myplanmodifier.UseStateWhen(func(a, b types.String) bool {
						return utils.NormalizeJson(a.ValueString()) == utils.NormalizeJson(b.ValueString())
					}),
				},
				DeprecationMessage: "This feature is deprecated and will be removed in a major release. Please use the `payload` argument to specify the body of the resource.",
			},

			"payload": schema.DynamicAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Dynamic{
					myplanmodifier.DynamicUseStateWhen(dynamic.SemanticallyEqual),
				},
			},

			"when": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.StringDefault("apply"),
				Validators: []validator.String{
					stringvalidator.OneOf("apply", "destroy"),
				},
			},

			"locks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
			},

			"response_export_values": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
			},

			"output": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "This feature is deprecated and will be removed in a major release. Please use the `output_payload` argument to output the response of the resource.",
			},

			"output_payload": schema.DynamicAttribute{
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

func (r *ActionResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var config *ActionResourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}
	// destroy doesn't need to modify plan
	if config == nil {
		return
	}

	// can't specify both body and payload
	if !config.Body.IsNull() && !config.Payload.IsNull() {
		response.Diagnostics.AddError("Invalid config", "can't specify both body and payload")
		return
	}
}

func (r *ActionResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	var config, plan, state *ActionResourceModel
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

	if state == nil || !plan.ResponseExportValues.Equal(state.ResponseExportValues) || utils.NormalizeJson(plan.Body.ValueString()) != utils.NormalizeJson(state.Body.ValueString()) || !plan.Payload.Equal(state.Payload) {
		plan.Output = types.StringUnknown()
		plan.OutputPayload = basetypes.NewDynamicUnknown()
		response.Diagnostics.Append(response.Plan.Set(ctx, plan)...)
	}
}

func (r *ActionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var model ActionResourceModel
	if response.Diagnostics.Append(request.Plan.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	if model.When.ValueString() == "apply" {
		r.Action(ctx, model, &response.State, &response.Diagnostics)
	} else {
		id, err := parse.ResourceIDWithResourceType(model.ResourceId.ValueString(), model.Type.ValueString())
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		resourceId := id.ID()
		if actionName := model.Action.ValueString(); actionName != "" {
			resourceId = fmt.Sprintf("%s/%s", id.ID(), actionName)
		}
		model.ID = basetypes.NewStringValue(resourceId)
		model.Output = basetypes.NewStringValue("{}")
		model.OutputPayload = basetypes.NewDynamicNull()
		response.Diagnostics.Append(response.State.Set(ctx, model)...)
	}
}

func (r *ActionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var model ActionResourceModel
	if response.Diagnostics.Append(request.Plan.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	if model.When.ValueString() == "apply" {
		r.Action(ctx, model, &response.State, &response.Diagnostics)
	}
}

func (r *ActionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var model ActionResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	if model.When.ValueString() != "destroy" {
		r.Action(ctx, model, &response.State, &response.Diagnostics)
	}
}

func (r *ActionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {

}

func (r *ActionResource) Action(ctx context.Context, model ActionResourceModel, state *tfsdk.State, diagnostics *diag.Diagnostics) {
	actionTimeout, diags := model.Timeouts.Create(ctx, 30*time.Minute)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, actionTimeout)
	defer cancel()

	id, err := parse.ResourceIDWithResourceType(model.ResourceId.ValueString(), model.Type.ValueString())
	if err != nil {
		diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	var requestBody interface{}
	switch {
	case !model.Payload.IsNull():
		out, err := expandPayload(model.Payload)
		if err != nil {
			diagnostics.AddError("Invalid payload", err.Error())
			return
		}
		requestBody = out
	case !model.Body.IsNull():
		bodyValueString := model.Body.ValueString()
		err := json.Unmarshal([]byte(bodyValueString), &requestBody)
		if err != nil {
			diagnostics.AddError("Invalid JSON string", fmt.Sprintf(`The argument "body" is invalid: value: %s, err: %+v`, model.Body.ValueString(), err))
			return
		}
	default:
		requestBody = map[string]interface{}{}
	}

	for _, id := range AsStringList(model.Locks) {
		locks.ByID(id)
		defer locks.UnlockByID(id)
	}

	client := r.ProviderData.ResourceClient
	responseBody, err := client.Action(ctx, id.AzureResourceId, model.Action.ValueString(), id.ApiVersion, model.Method.ValueString(), requestBody)
	if err != nil {
		diagnostics.AddError("Failed to perform action", fmt.Errorf("performing action %s of %q: %+v", model.Action.ValueString(), id, err).Error())
		return
	}

	resourceId := id.ID()
	if actionName := model.Action.ValueString(); actionName != "" {
		resourceId = fmt.Sprintf("%s/%s", id.ID(), actionName)
	}
	model.ID = basetypes.NewStringValue(resourceId)
	model.Output = basetypes.NewStringValue(flattenOutput(responseBody, AsStringList(model.ResponseExportValues)))
	model.OutputPayload = types.DynamicValue(flattenOutputPayload(responseBody, AsStringList(model.ResponseExportValues)))

	diagnostics.Append(state.Set(ctx, model)...)
}
