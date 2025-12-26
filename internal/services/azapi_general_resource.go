package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/Azure/terraform-provider-azapi/internal/services/defaults"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier/planmodifierdynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/skip"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	urlSchemeDivider = "://"
)

type GeneralResourceModel struct {
	ID                            types.String     `tfsdk:"id"`
	Body                          types.Dynamic    `tfsdk:"body"`
	SensitiveBody                 types.Dynamic    `tfsdk:"sensitive_body"`
	SensitiveBodyVersion          types.Map        `tfsdk:"sensitive_body_version"`
	Url                           types.String     `tfsdk:"url"`
	ApiVersion                    types.String     `tfsdk:"api_version"`
	AuthEndpoint                  types.String     `tfsdk:"auth_endpoint"`
	AuthAudience                  types.String     `tfsdk:"auth_audience"`
	IgnoreCasing                  types.Bool       `tfsdk:"ignore_casing"`
	IgnoreMissingProperty         types.Bool       `tfsdk:"ignore_missing_property"`
	IgnoreNullProperty            types.Bool       `tfsdk:"ignore_null_property"`
	ReplaceTriggersExternalValues types.Dynamic    `tfsdk:"replace_triggers_external_values"`
	ReplaceTriggersRefs           types.List       `tfsdk:"replace_triggers_refs"`
	ResponseExportValues          types.Dynamic    `tfsdk:"response_export_values"`
	Retry                         retry.RetryValue `tfsdk:"retry" skip_on:"update"`
	Locks                         types.List       `tfsdk:"locks"`
	Output                        types.Dynamic    `tfsdk:"output"`
	Timeouts                      timeouts.Value   `tfsdk:"timeouts" skip_on:"update"`
	CreateHeaders                 types.Map        `tfsdk:"create_headers" skip_on:"update"`
	CreateQueryParameters         types.Map        `tfsdk:"create_query_parameters" skip_on:"update"`
	UpdateHeaders                 types.Map        `tfsdk:"update_headers"`
	UpdateQueryParameters         types.Map        `tfsdk:"update_query_parameters"`
	DeleteHeaders                 types.Map        `tfsdk:"delete_headers" skip_on:"update"`
	DeleteQueryParameters         types.Map        `tfsdk:"delete_query_parameters" skip_on:"update"`
	ReadHeaders                   types.Map        `tfsdk:"read_headers" skip_on:"update"`
	ReadQueryParameters           types.Map        `tfsdk:"read_query_parameters" skip_on:"update"`
}

type GeneralResource struct {
	ProviderData *clients.Client
}

var _ resource.Resource = &GeneralResource{}
var _ resource.ResourceWithConfigure = &GeneralResource{}
var _ resource.ResourceWithModifyPlan = &GeneralResource{}
var _ resource.ResourceWithImportState = &GeneralResource{}
var _ resource.ResourceWithMoveState = &GeneralResource{}

func (r *GeneralResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring azapi_general_resource")
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *GeneralResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_general_resource"
}

func (r *GeneralResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This resource can manage general resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: docstrings.ID(),
			},

			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The URL of the data plane resource. Must include `https://`, e.g. `https://my-keyvault.vault.azure.net/secrets/test-secret/b15f7ed819404dd0b28debe0aa711a54`.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^https://.*`), "The URL must start with `https://`. It should be a relative URL, e.g. `https://my-keyvault.vault.azure.net/secrets/test-secret/b15f7ed819404dd0b28debe0aa711a54`."),
				},
			},

			"api_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The API version to use for the data plane resource. If not specified, the provider will not add the api-version to the request, which may be required.",
			},

			"auth_endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional authentication endpoint to use for the data plane resource. If not specified, the provider will use the base URL.",
			},

			"auth_audience": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional authentication audience to use for the data plane resource. If not specified, the provider will use the base URL.",
			},

			// The body attribute is a dynamic attribute that only allows users to specify the resource body as an HCL object
			"body": schema.DynamicAttribute{
				Optional: true,
				Computed: true,
				// in the previous version, the default value is string "{}", now it's a dynamic value {}
				Default: defaults.DynamicDefault(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{})),
				PlanModifiers: []planmodifier.Dynamic{
					myplanmodifier.DynamicUseStateWhen(dynamic.SemanticallyEqual),
				},
				MarkdownDescription: docstrings.Body(),
				Validators: []validator.Dynamic{
					myvalidator.DynamicIsNotStringValidator(),
				},
			},

			"sensitive_body": schema.DynamicAttribute{
				Optional:            true,
				WriteOnly:           true,
				MarkdownDescription: docstrings.SensitiveBody(),
			},

			"sensitive_body_version": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: docstrings.SensitiveBodyVersion(),
			},

			"ignore_casing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.BoolDefault(false),
				MarkdownDescription: docstrings.Body(),
			},

			"ignore_missing_property": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.BoolDefault(true),
				MarkdownDescription: docstrings.IgnoreMissingProperty(),
			},

			"ignore_null_property": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.BoolDefault(false),
				MarkdownDescription: docstrings.IgnoreNullProperty(),
			},

			"response_export_values": schema.DynamicAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Dynamic{
					myplanmodifier.DynamicUseStateWhen(dynamic.SemanticallyEqual),
				},
				MarkdownDescription: docstrings.ResponseExportValues(),
			},

			"retry": retry.RetrySchema(ctx),

			"replace_triggers_external_values": schema.DynamicAttribute{
				Optional: true,
				MarkdownDescription: "Will trigger a replace of the resource when the value changes and is not `null`. This can be used by practitioners to force a replace of the resource when certain values change, e.g. changing the SKU of a virtual machine based on the value of variables or locals. " +
					"The value is a `dynamic`, so practitioners can compose the input however they wish. For a \"break glass\" set the value to `null` to prevent the plan modifier taking effect. \n" +
					"If you have `null` values that you do want to be tracked as affecting the resource replacement, include these inside an object. \n" +
					"Advanced use cases are possible and resource replacement can be triggered by values external to the resource, for example when a dependent resource changes.\n\n" +
					"e.g. to replace a resource when either the SKU or os_type attributes change:\n" +
					"\n" +
					"```hcl\n" +
					"resource \"azapi_data_plane_resource\" \"example\" {\n" +
					"  name = var.name\n" +
					"  type = \"Microsoft.AppConfiguration/configurationStores/keyValues@1.0\"\n" +
					"  body = {\n" +
					"    properties = {\n" +
					"      sku   = var.sku\n" +
					"      zones = var.zones\n" +
					"    }\n" +
					"  }\n" +
					"\n" +
					"  replace_triggers_external_values = [\n" +
					"    var.sku,\n" +
					"    var.zones,\n" +
					"  ]\n" +
					"}\n" +
					"```\n",
				PlanModifiers: []planmodifier.Dynamic{
					planmodifierdynamic.RequiresReplaceIfNotNull(),
				},
			},

			"replace_triggers_refs": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A list of paths in the current Terraform configuration. When the values at these paths change, the resource will be replaced.",
			},

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
				MarkdownDescription: docstrings.Output("azapi_data_plane_resource"),
			},

			"create_headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A mapping of headers to be sent with the create request.",
			},

			"create_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the create request.",
			},

			"update_headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A mapping of headers to be sent with the update request.",
			},

			"update_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the update request.",
			},

			"delete_headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A mapping of headers to be sent with the delete request.",
			},

			"delete_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the delete request.",
			},

			"read_headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A mapping of headers to be sent with the read request.",
			},

			"read_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the read request.",
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

func (r *GeneralResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	var config, plan, state *GeneralResourceModel
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

	// In the below two cases, we think the config is still matched with the remote state, and there's no need to update the resource:
	// 1. If the api-version is changed, but the body is not changed
	// 2. If the body only removes/adds properties that are equal to the remote state
	if r.ProviderData.Features.IgnoreNoOpChanges && dynamic.IsFullyKnown(plan.Body) && state != nil && (!dynamic.SemanticallyEqual(plan.Body, state.Body) || !plan.ApiVersion.Equal(state.ApiVersion)) {
		// GET the existing resource with config's api-version
		serviceConfig, err := generateServiceConfig(state)
		if err != nil {
			response.Diagnostics.AddError("Failed to generate service configuration", err.Error())
			return
		}
		responseBody, err := r.ProviderData.GeneralClient.Get(ctx, state.Url.ValueString(), state.ApiVersion.ValueString(), serviceConfig, clients.DefaultRequestOptions())
		if err != nil {
			response.Diagnostics.AddError("Failed to retrieve resource", fmt.Sprintf("Retrieving existing resource %s: %+v", state.ID.ValueString(), err))
			return
		}
		stateBody := make(map[string]interface{})
		if err := unmarshalBody(state.Body, &stateBody); err != nil {
			response.Diagnostics.AddError("Invalid state body", fmt.Sprintf(`The argument "body" in state is invalid: %s`, err.Error()))
			return
		}
		// stateBody contains sensitive properties that are not returned in GET response
		responseBody = utils.MergeObject(responseBody, stateBody)

		configBody := make(map[string]interface{})
		if err := unmarshalBody(plan.Body, &configBody); err != nil {
			response.Diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
			return
		}
		option := utils.UpdateJsonOption{
			IgnoreCasing:          plan.IgnoreCasing.ValueBool(),
			IgnoreMissingProperty: false,
			IgnoreNullProperty:    plan.IgnoreNullProperty.ValueBool(),
		}
		remoteBody := utils.UpdateObject(configBody, responseBody, option)
		// suppress the change if the remote body is equal to the config body
		if reflect.DeepEqual(remoteBody, configBody) {
			plan.Body = state.Body
			plan.ApiVersion = state.ApiVersion
		}
	}

	if state != nil {
		// Set output as unknown to trigger a plan diff, if ephemeral body has changed
		diff, diags := ephemeralBodyChangeInPlan(ctx, request.Private, config.SensitiveBody, config.SensitiveBodyVersion, state.SensitiveBodyVersion)
		if response.Diagnostics = append(response.Diagnostics, diags...); response.Diagnostics.HasError() {
			return
		}
		if diff {
			tflog.Info(ctx, `"sensitive_body" has changed`)
			plan.Output = types.DynamicUnknown()
		}
	}

	// Check if any paths in replace_triggers_refs have changed
	if state != nil && plan != nil && !plan.ReplaceTriggersRefs.IsNull() {
		refPaths := make(map[string]string)
		for pathIndex, refPath := range common.AsStringList(plan.ReplaceTriggersRefs) {
			refPaths[fmt.Sprintf("%d", pathIndex)] = refPath
		}

		// read previous values from state
		stateData, err := dynamic.ToJSON(state.Body)
		if err != nil {
			response.Diagnostics.AddError("Invalid state body configuration", err.Error())
			return
		}
		var stateModel interface{}
		err = json.Unmarshal(stateData, &stateModel)
		if err != nil {
			response.Diagnostics.AddError("Invalid state body configuration", err.Error())
			return
		}
		previousValues := flattenOutputJMES(stateModel, refPaths)

		// read current values from plan
		planData, err := dynamic.ToJSON(plan.Body)
		if err != nil {
			response.Diagnostics.AddError("Invalid plan body configuration", err.Error())
			return
		}
		var planModel interface{}
		err = json.Unmarshal(planData, &planModel)
		if err != nil {
			response.Diagnostics.AddError("Invalid plan body configuration", err.Error())
			return
		}
		currentValues := flattenOutputJMES(planModel, refPaths)

		// compare previous and current values
		if !reflect.DeepEqual(previousValues, currentValues) {
			response.RequiresReplace.Append(path.Root("body"))
		}
	}
}

func (r *GeneralResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	r.CreateUpdate(ctx, request.Config, request.Plan, &response.State, &response.Diagnostics)
}

func (r *GeneralResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// See if we can skip the external API call (changes are to state only)
	var state, plan GeneralResourceModel
	if response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...); response.Diagnostics.HasError() {
		return
	}
	if response.Diagnostics.Append(request.State.Get(ctx, &state)...); response.Diagnostics.HasError() {
		return
	}
	if skip.CanSkipExternalRequest(state, plan, "update") {
		tflog.Debug(ctx, "azapi_resource.CreateUpdate skipping external request as no unskippable changes were detected")
		response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	}
	tflog.Debug(ctx, "azapi_resource.CreateUpdate proceeding with external request as no skippable changes were detected")
	r.CreateUpdate(ctx, request.Config, request.Plan, &response.State, &response.Diagnostics)
}

func (r *GeneralResource) CreateUpdate(ctx context.Context, requestConfig tfsdk.Config, plan tfsdk.Plan, state *tfsdk.State, diagnostics *diag.Diagnostics) {
	var planModel, stateModel, configModel GeneralResourceModel
	diagnostics.Append(requestConfig.Get(ctx, &configModel)...)
	diagnostics.Append(plan.Get(ctx, &planModel)...)
	diagnostics.Append(state.Get(ctx, &stateModel)...)
	if diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "url", planModel.Url.ValueString())

	urlParsed, err := url.Parse(planModel.Url.ValueString())
	if err != nil {
		diagnostics.AddError("Invalid URL", err.Error())
		return
	}

	isNewResource := state == nil || state.Raw.IsNull()

	var timeout time.Duration
	var diags diag.Diagnostics
	if isNewResource {
		timeout, diags = planModel.Timeouts.Create(ctx, 30*time.Minute)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return
		}
	} else {
		timeout, diags = planModel.Timeouts.Update(ctx, 30*time.Minute)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return
		}
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Ensure the context deadline has been set before calling ConfigureClientWithCustomRetry().
	client := r.ProviderData.GeneralClient.ConfigureClientWithCustomRetry(ctx, planModel.Retry, false)

	serviceConfig, err := generateServiceConfig(&planModel)
	if err != nil {
		diagnostics.AddError("Failed to generate service configuration", err.Error())
		return
	}

	if isNewResource {
		// check if the resource already exists using the non-retry client to avoid issue where user specifies
		// a FooResourceNotFound error as a retryable error
		_, err = r.ProviderData.GeneralClient.Get(ctx, urlParsed.String(), planModel.ApiVersion.ValueString(), serviceConfig, clients.NewRequestOptions(common.AsMapOfString(planModel.ReadHeaders), common.AsMapOfLists(planModel.ReadQueryParameters)))
		if err == nil {
			diagnostics.AddError("Resource already exists", tf.ImportAsExistsError("azapi_general_resource", urlParsed.String()).Error())
			return
		}
		if !utils.ResponseErrorWasNotFound(err) {
			diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("checking for presence of existing %s: %+v", urlParsed.String(), err).Error())
			return
		}
	}

	// build the request body
	body := make(map[string]interface{})
	if err := unmarshalBody(planModel.Body, &body); err != nil {
		diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
		return
	}
	sensitiveBodyVersionInState := types.MapNull(types.StringType)
	if state != nil {
		sensitiveBodyVersionInState = stateModel.SensitiveBodyVersion
	}
	sensitiveBody, err := unmarshalSensitiveBody(configModel.SensitiveBody, planModel.SensitiveBodyVersion, sensitiveBodyVersionInState)
	if err != nil {
		diagnostics.AddError("Invalid sensitive_body", fmt.Sprintf(`The argument "sensitive_body" is invalid: %s`, err.Error()))
		return
	}

	if sensitiveBody != nil {
		body = utils.MergeObject(body, sensitiveBody).(map[string]interface{})
	}

	if planModel.IgnoreNullProperty.ValueBool() {
		out := utils.RemoveNullProperty(body)
		v, ok := out.(map[string]interface{})
		if ok {
			body = v
		}
	}

	// create/update the resource
	lockIds := common.AsStringList(planModel.Locks)
	slices.Sort(lockIds)
	for _, lockId := range lockIds {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	options := clients.NewRequestOptions(common.AsMapOfString(planModel.CreateHeaders), common.AsMapOfLists(planModel.CreateQueryParameters))
	_, err = client.CreateOrUpdateThenPoll(ctx, urlParsed.String(), planModel.ApiVersion.ValueString(), serviceConfig, body, options)
	if err != nil {
		diagnostics.AddError("Failed to create/update resource", fmt.Errorf("creating/updating %q: %+v", urlParsed.String(), err).Error())
		return
	}

	// Create a new retry client to handle specific case of transient 403/404 after resource creation
	// If a read after create retry is not specified, use the default.
	clientGetAfterPut := r.ProviderData.GeneralClient.ConfigureClientWithCustomRetry(ctx, planModel.Retry, true)

	responseBody, err := clientGetAfterPut.Get(ctx, urlParsed.String(), planModel.ApiVersion.ValueString(), serviceConfig, clients.NewRequestOptions(common.AsMapOfString(planModel.ReadHeaders), common.AsMapOfLists(planModel.ReadQueryParameters)))
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("Error reading %q - removing from state", urlParsed.String()))
			state.RemoveResource(ctx)
			return
		}
		diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", urlParsed.String(), err).Error())
		return
	}
	planModel.ID = basetypes.NewStringValue(urlParsed.String() + "@" + planModel.ApiVersion.ValueString())

	output, err := buildOutputFromBody(responseBody, planModel.ResponseExportValues, nil)
	if err != nil {
		diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	planModel.Output = output

	diagnostics.Append(state.Set(ctx, planModel)...)
}

func (r *GeneralResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var model GeneralResourceModel
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

	ctx = tflog.SetField(ctx, "url", model.Url.ValueString())

	urlParsed, err := url.Parse(model.Url.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid URL", err.Error())
		return
	}

	// Ensure that the context deadline has been set before calling ConfigureClientWithCustomRetry().
	client := r.ProviderData.GeneralClient.ConfigureClientWithCustomRetry(ctx, model.Retry, false)

	serviceConfig, err := generateServiceConfig(&model)
	if err != nil {
		response.Diagnostics.AddError("Failed to generate service configuration", err.Error())
		return
	}

	responseBody, err := client.Get(ctx, urlParsed.String(), model.ApiVersion.ValueString(), serviceConfig, clients.NewRequestOptions(common.AsMapOfString(model.ReadHeaders), common.AsMapOfLists(model.ReadQueryParameters)))
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("[INFO] Error reading %q - removing from state", urlParsed.String()))
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", urlParsed.String(), err).Error())
		return
	}

	requestBody := make(map[string]interface{})
	if err := unmarshalBody(model.Body, &requestBody); err != nil {
		response.Diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
		return
	}

	option := utils.UpdateJsonOption{
		IgnoreCasing:          model.IgnoreCasing.ValueBool(),
		IgnoreMissingProperty: model.IgnoreMissingProperty.ValueBool(),
	}
	body := utils.UpdateObject(requestBody, responseBody, option)

	data, err := json.Marshal(body)
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

	if !model.Body.IsNull() {
		payload, err := dynamic.FromJSON(data, model.Body.UnderlyingValue().Type(ctx))
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to parse payload: %s", err.Error()))
			payload, err = dynamic.FromJSONImplied(data)
			if err != nil {
				response.Diagnostics.AddError("Invalid payload", err.Error())
				return
			}
		}
		model.Body = payload
	}

	response.Diagnostics.Append(response.State.Set(ctx, model)...)
}

func (r *GeneralResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var model GeneralResourceModel
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

	urlParsed, err := url.Parse(model.Url.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid URL", err.Error())
		return
	}

	ctx = tflog.SetField(ctx, "resource_id", urlParsed.String())

	// Ensure the context deadline has been set before calling ConfigureClientWithCustomRetry().
	client := r.ProviderData.GeneralClient.ConfigureClientWithCustomRetry(ctx, model.Retry, false)

	serviceConfig, err := generateServiceConfig(&model)
	if err != nil {
		response.Diagnostics.AddError("Failed to generate service configuration", err.Error())
		return
	}

	lockIds := common.AsStringList(model.Locks)
	slices.Sort(lockIds)
	for _, lockId := range lockIds {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	_, err = client.DeleteThenPoll(ctx, urlParsed.String(), model.ApiVersion.ValueString(), serviceConfig, clients.NewRequestOptions(common.AsMapOfString(model.DeleteHeaders), common.AsMapOfLists(model.DeleteQueryParameters)))
	if err != nil && !utils.ResponseErrorWasNotFound(err) {
		response.Diagnostics.AddError("Failed to delete resource", fmt.Errorf("deleting %s: %+v", urlParsed.String(), err).Error())
	}
}

func (r *GeneralResource) MoveState(ctx context.Context) []resource.StateMover {
	return []resource.StateMover{
		{
			SourceSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			StateMover: func(ctx context.Context, request resource.MoveStateRequest, response *resource.MoveStateResponse) {
				if !strings.HasPrefix(request.SourceTypeName, "azurerm") {
					response.Diagnostics.AddError("Invalid source type", "The `azapi_general_resource` resource can only be moved from an `azurerm` resource")
					return
				}

				if request.SourceState == nil {
					response.Diagnostics.AddError("Invalid source state", "The source state is nil")
					return
				}

				requestID := ""
				if response.Diagnostics.Append(request.SourceState.GetAttribute(ctx, path.Root("id"), &requestID)...); response.Diagnostics.HasError() {
					return
				}
				if requestID == "" {
					response.Diagnostics.AddError("Invalid source state", "The source state does not contain an id")
					return
				}

				state := r.defaultAzapiResourceModel()
				state.ID = types.StringValue(requestID)

				response.Diagnostics.Append(response.TargetPrivate.SetKey(ctx, FlagMoveState, []byte("true"))...)
				response.Diagnostics.Append(response.TargetState.Set(ctx, state)...)
			},
		},
	}
}

func (r *GeneralResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Importing Resource - parsing %q", request.ID))

	client := r.ProviderData.GeneralClient

	state := r.defaultAzapiResourceModel()
	state.ID = types.StringValue(request.ID)
	url, apiVersion, _ := strings.Cut(request.ID, "@")
	state.Url = types.StringValue(url)
	if apiVersion != "" {
		state.ApiVersion = types.StringValue(apiVersion)
	}

	svcConfig, err := generateServiceConfig(&state)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to generate service configuration from url %q and api version %q", url, apiVersion), err.Error())
		return
	}

	responseBody, err := client.Get(ctx, url, apiVersion, svcConfig, clients.NewRequestOptions(common.AsMapOfString(state.ReadHeaders), common.AsMapOfLists(state.ReadQueryParameters)))
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("[INFO] Error reading %q - removing from state", url))
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", url, err).Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("resource %q is imported", url))

	respondeJson, err := json.Marshal(responseBody)
	if err != nil {
		response.Diagnostics.AddError("Failed to marshal response body", err.Error())
		return
	}

	payload, err := dynamic.FromJSONImplied(respondeJson)
	if err != nil {
		response.Diagnostics.AddError("Invalid body", err.Error())
		return
	}
	state.Body = payload

	var defaultOutput interface{}
	output, err := buildOutputFromBody(responseBody, state.ResponseExportValues, defaultOutput)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	state.Output = output

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

// generateServiceConfig generates the service configuration for the GeneralResourceModel.
// It is used in the client's cache to store the HTTP pipeline configuration for a given authentication
// endpoint and audience.
func generateServiceConfig(model *GeneralResourceModel) (cloud.ServiceConfiguration, error) {
	urlParsed, err := url.Parse(model.Url.ValueString())
	if err != nil {
		return cloud.ServiceConfiguration{}, err
	}
	authAudience := model.AuthAudience.ValueString()
	if authAudience == "" {
		authAudience = urlParsed.Scheme + urlSchemeDivider + urlParsed.Hostname()
	}
	authEndpoint := model.AuthEndpoint.ValueString()
	if authEndpoint == "" {
		authEndpoint = urlParsed.Scheme + urlSchemeDivider + urlParsed.Hostname()
	}

	return cloud.ServiceConfiguration{
		Endpoint: authEndpoint,
		Audience: authAudience,
	}, nil
}

func (r *GeneralResource) defaultAzapiResourceModel() GeneralResourceModel {
	return GeneralResourceModel{
		ID:                            types.StringNull(),
		Body:                          types.Dynamic{},
		SensitiveBody:                 types.DynamicNull(),
		SensitiveBodyVersion:          types.MapNull(types.StringType),
		IgnoreCasing:                  types.BoolValue(false),
		IgnoreMissingProperty:         types.BoolValue(true),
		IgnoreNullProperty:            types.BoolValue(false),
		Locks:                         types.ListNull(types.StringType),
		Output:                        types.DynamicNull(),
		ReplaceTriggersExternalValues: types.DynamicNull(),
		ReplaceTriggersRefs:           types.ListNull(types.StringType),
		ResponseExportValues:          types.DynamicNull(),
		Retry:                         retry.RetryValue{},
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"read":   types.StringType,
				"delete": types.StringType,
			}),
		},
		CreateHeaders:         types.MapNull(types.StringType),
		CreateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		UpdateHeaders:         types.MapNull(types.StringType),
		UpdateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		DeleteHeaders:         types.MapNull(types.StringType),
		DeleteQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		ReadHeaders:           types.MapNull(types.StringType),
		ReadQueryParameters:   types.MapNull(types.ListType{ElemType: types.StringType}),
	}
}
