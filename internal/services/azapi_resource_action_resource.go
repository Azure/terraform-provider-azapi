package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/defaults"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/migration"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
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
	ID                   types.String        `tfsdk:"id"`
	Type                 types.String        `tfsdk:"type"`
	ResourceId           types.String        `tfsdk:"resource_id"`
	Action               types.String        `tfsdk:"action"`
	Method               types.String        `tfsdk:"method"`
	Body                 types.Dynamic       `tfsdk:"body"`
	When                 types.String        `tfsdk:"when"`
	Locks                types.List          `tfsdk:"locks"`
	ResponseExportValues types.Dynamic       `tfsdk:"response_export_values"`
	Output               types.Dynamic       `tfsdk:"output"`
	Timeouts             timeouts.Value      `tfsdk:"timeouts"`
	Retry                retry.RetryValue    `tfsdk:"retry"`
	Headers              map[string]string   `tfsdk:"headers"`
	QueryParameters      map[string][]string `tfsdk:"query_parameters"`
}

type ActionResource struct {
	ProviderData *clients.Client
}

var _ resource.Resource = &ActionResource{}
var _ resource.ResourceWithConfigure = &ActionResource{}
var _ resource.ResourceWithModifyPlan = &ActionResource{}
var _ resource.ResourceWithUpgradeState = &ActionResource{}

func (r *ActionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *ActionResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_action"
}

func (r *ActionResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: migration.AzapiResourceActionMigrationV0ToV2(ctx),
		1: migration.AzapiResourceActionMigrationV1ToV2(ctx),
	}
}

func (r *ActionResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: docstrings.ID(),
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: docstrings.Type(),
			},

			"resource_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: "The ID of an existing Azure source.",
			},

			"action": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: docstrings.ResourceAction(),
			},

			"method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.StringDefault("POST"),
				Validators: []validator.String{
					stringvalidator.OneOf("POST", "PATCH", "PUT", "DELETE", "GET", "HEAD"),
				},
				MarkdownDescription: "Specifies the HTTP method of the azure resource action. Allowed values are `POST`, `PATCH`, `PUT` and `DELETE`. Defaults to `POST`.",
			},

			// The body attribute is a dynamic attribute that only allows users to specify the resource body as an HCL object
			"body": schema.DynamicAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Dynamic{
					myplanmodifier.DynamicUseStateWhen(dynamic.SemanticallyEqual),
				},
				MarkdownDescription: docstrings.Body(),
				Validators: []validator.Dynamic{
					myvalidator.DynamicIsNotStringValidator(),
				},
			},

			"when": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.StringDefault("apply"),
				Validators: []validator.String{
					stringvalidator.OneOf("apply", "destroy"),
				},
				MarkdownDescription: "When to perform the action, value must be one of: `apply`, `destroy`. Default is `apply`.",
			},

			"locks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
				MarkdownDescription: docstrings.Locks(),
			},

			"response_export_values": CommonAttributeResponseExportValues(),

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("azapi_resource_action"),
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
				Create: true,
				Update: true,
				Read:   true,
				Delete: true,
			}),
		},

		Version: 2,
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

	if state == nil || !plan.ResponseExportValues.Equal(state.ResponseExportValues) || !dynamic.SemanticallyEqual(plan.Body, state.Body) {
		plan.Output = basetypes.NewDynamicUnknown()
	} else {
		plan.Output = state.Output
	}

	response.Diagnostics.Append(response.Plan.Set(ctx, plan)...)
}

func (r *ActionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var model ActionResourceModel
	if response.Diagnostics.Append(request.Plan.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	timeout, diags := model.Timeouts.Create(ctx, 30*time.Minute)
	if response.Diagnostics.Append(diags...); response.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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
		model.Output = basetypes.NewDynamicNull()
		response.Diagnostics.Append(response.State.Set(ctx, model)...)
	}
}

func (r *ActionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var model ActionResourceModel
	if response.Diagnostics.Append(request.Plan.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	timeout, diags := model.Timeouts.Update(ctx, 30*time.Minute)
	if response.Diagnostics.Append(diags...); response.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if model.When.ValueString() == "apply" {
		r.Action(ctx, model, &response.State, &response.Diagnostics)
	}
}

func (r *ActionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var model ActionResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	if model.When.ValueString() == "destroy" {
		r.Action(ctx, model, &response.State, &response.Diagnostics)
	}
}

func (r *ActionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state ActionResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &state)...); response.Diagnostics.HasError() {
		return
	}

	if state.When.IsNull() {
		state.When = basetypes.NewStringValue("apply")
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
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
	if err := unmarshalBody(model.Body, &requestBody); err != nil {
		diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
		return
	}

	for _, id := range AsStringList(model.Locks) {
		locks.ByID(id)
		defer locks.UnlockByID(id)
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

	responseBody, err := client.Action(ctx, id.AzureResourceId, model.Action.ValueString(), id.ApiVersion, model.Method.ValueString(), requestBody, clients.NewRequestOptions(model.Headers, model.QueryParameters))
	if err != nil {
		diagnostics.AddError("Failed to perform action", fmt.Errorf("performing action %s of %q: %+v", model.Action.ValueString(), id, err).Error())
		return
	}

	resourceId := id.ID()
	if actionName := model.Action.ValueString(); actionName != "" {
		resourceId = fmt.Sprintf("%s/%s", id.ID(), actionName)
	}
	model.ID = basetypes.NewStringValue(resourceId)

	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, nil)
	if err != nil {
		diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	diagnostics.Append(state.Set(ctx, model)...)
}
