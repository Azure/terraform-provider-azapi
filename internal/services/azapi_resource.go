package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/services/defaults"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

type AzapiResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ParentID                types.String `tfsdk:"parent_id"`
	Type                    types.String `tfsdk:"type"`
	Location                types.String `tfsdk:"location"`
	Identity                types.List   `tfsdk:"identity"`
	Body                    types.String `tfsdk:"body"`
	Locks                   types.List   `tfsdk:"locks"`
	RemovingSpecialChars    types.Bool   `tfsdk:"removing_special_chars"`
	SchemaValidationEnabled types.Bool   `tfsdk:"schema_validation_enabled"`
	IgnoreBodyChanges       types.List   `tfsdk:"ignore_body_changes"`
	IgnoreCasing            types.Bool   `tfsdk:"ignore_casing"`
	IgnoreMissingProperty   types.Bool   `tfsdk:"ignore_missing_property"`
	ResponseExportValues    types.List   `tfsdk:"response_export_values"`
	Output                  types.String `tfsdk:"output"`
	Tags                    types.Map    `tfsdk:"tags"`
}

var _ resource.Resource = &AzapiResource{}
var _ resource.ResourceWithConfigure = &AzapiResource{}
var _ resource.ResourceWithModifyPlan = &AzapiResource{}
var _ resource.ResourceWithValidateConfig = &AzapiResource{}
var _ resource.ResourceWithImportState = &AzapiResource{}

type AzapiResource struct {
	ProviderData *clients.Client
}

func (r *AzapiResource) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *AzapiResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource"
}

func (r *AzapiResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"removing_special_chars": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.BoolDefault(false),
			},

			"parent_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
			},

			"location": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					myplanmodifier.UseStateWhen(func(a, b types.String) bool {
						return location.Normalize(a.ValueString()) == location.Normalize(b.ValueString())
					}),
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

			"ignore_body_changes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
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

			"schema_validation_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  defaults.BoolDefault(true),
			},

			"output": schema.StringAttribute{
				Computed: true,
			},

			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.Map{
					tags.Validator(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"identity": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{stringvalidator.OneOf(
								string(identity.SystemAssignedUserAssigned),
								string(identity.UserAssigned),
								string(identity.SystemAssigned),
								string(identity.None),
							)},
						},

						"identity_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(myvalidator.StringIsUserAssignedIdentityID()),
							},
						},

						"principal_id": schema.StringAttribute{
							Computed: true,
						},

						"tenant_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
		Version: 0,
	}
}

func (r *AzapiResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var config *AzapiResourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}
	// destroy doesn't need to modify plan
	if config == nil {
		return
	}

	resourceType := config.Type.ValueString()

	// for resource group, if parent_id is not specified, set it to subscription id
	if config.ParentID.IsNull() {
		azureResourceType, _, _ := utils.GetAzureResourceTypeApiVersion(resourceType)
		if !strings.EqualFold(azureResourceType, arm.ResourceGroupResourceType.String()) {
			response.Diagnostics.AddError("Missing required argument", `The argument "parent_id" is required, but no definition was found.`)
			return
		}
	}

	if config.Body.IsUnknown() {
		return
	}

	body := make(map[string]interface{})
	if bodyValueString := config.Body.ValueString(); bodyValueString != "" {
		if err := json.Unmarshal([]byte(bodyValueString), &body); err != nil {
			response.Diagnostics.AddError("Invalid JSON string", fmt.Sprintf(`The argument "body" is invalid: value: %s, err: %+v`, bodyValueString, err))
			return
		}
	}
	if diags := validateDuplicatedDefinitions(config, body); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

func (r *AzapiResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	var config, state, plan *AzapiResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// destroy doesn't need to modify plan
	if config == nil {
		return
	}

	defer func() {
		response.Plan.Set(ctx, plan)
	}()

	resourceType := config.Type.ValueString()

	// for resource group, if parent_id is not specified, set it to subscription id
	if config.ParentID.IsNull() {
		azureResourceType, _, _ := utils.GetAzureResourceTypeApiVersion(resourceType)
		if strings.EqualFold(azureResourceType, arm.ResourceGroupResourceType.String()) {
			plan.ParentID = types.StringValue(fmt.Sprintf("/subscriptions/%s", r.ProviderData.Account.GetSubscriptionId()))
		}
	}

	if assignedName, diags := r.nameWithDefaultNaming(config.Name, plan.RemovingSpecialChars.ValueBool()); !diags.HasError() {
		plan.Name = assignedName
		// replace the resource if the name is changed
		if state != nil && !state.Name.Equal(plan.Name) {
			response.RequiresReplace.Append(path.Root("name"))
		}
	} else {
		response.Diagnostics.Append(diags...)
		return
	}

	// if the config identity type and identity ids are not changed, use the state identity
	if !config.Identity.IsNull() && state != nil && !state.Identity.IsNull() {
		configIdentity := identity.FromList(config.Identity)
		stateIdentity := identity.FromList(state.Identity)
		if configIdentity.Type.Equal(stateIdentity.Type) && configIdentity.IdentityIDs.Equal(stateIdentity.IdentityIDs) {
			plan.Identity = state.Identity
		}
	}

	if plan.Body.IsUnknown() {
		if config.Tags.IsNull() {
			plan.Tags = basetypes.NewMapUnknown(types.StringType)
		}
		if config.Location.IsNull() {
			plan.Location = basetypes.NewStringUnknown()
		}
		plan.Output = types.StringUnknown()
		return
	}

	if state == nil || !plan.Identity.Equal(state.Identity) || !plan.ResponseExportValues.Equal(state.ResponseExportValues) || utils.NormalizeJson(plan.Body.ValueString()) != utils.NormalizeJson(state.Body.ValueString()) {
		plan.Output = types.StringUnknown()
	}

	body := make(map[string]interface{})
	err := json.Unmarshal([]byte(plan.Body.ValueString()), &body)
	if err != nil {
		response.Diagnostics.AddError("Invalid JSON string", fmt.Sprintf(`The argument "body" is invalid: value: %s, err: %+v`, plan.Body.ValueString(), err))
		return
	}

	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(config.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", fmt.Sprintf(`The argument "type" is invalid: %s`, err.Error()))
		return
	}
	resourceDef, _ := azure.GetResourceDefinition(azureResourceType, apiVersion)

	plan.Tags = r.tagsWithDefaultTags(config.Tags, body, state, resourceDef)
	if state == nil || !state.Tags.Equal(plan.Tags) {
		plan.Output = types.StringUnknown()
	}

	locationValue := plan.Location
	// the resource is new
	if locationValue.IsUnknown() || config.Location.IsNull() {
		locationValue = config.Location
	}
	plan.Location = r.locationWithDefaultLocation(locationValue, body, state, resourceDef)
	if state != nil && location.Normalize(state.Location.ValueString()) != location.Normalize(plan.Location.ValueString()) {
		// if the location is changed, replace the resource
		response.RequiresReplace.Append(path.Root("location"))
	}

	if plan.SchemaValidationEnabled.ValueBool() {
		if response.Diagnostics.Append(expandBody(body, *plan)...); response.Diagnostics.HasError() {
			return
		}
		body["name"] = plan.Name.ValueString()
		err = schemaValidation(azureResourceType, apiVersion, resourceDef, body)
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
	}
}

func (r *AzapiResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	r.CreateUpdate(ctx, request.Plan, &response.State, &response.Diagnostics)
}

func (r *AzapiResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	r.CreateUpdate(ctx, request.Plan, &response.State, &response.Diagnostics)
}

func (r *AzapiResource) CreateUpdate(ctx context.Context, requestPlan tfsdk.Plan, responseState *tfsdk.State, diagnostics *diag.Diagnostics) {
	var plan, state *AzapiResourceModel
	diagnostics.Append(requestPlan.Get(ctx, &plan)...)
	diagnostics.Append(responseState.Get(ctx, &state)...)
	if diagnostics.HasError() {
		return
	}

	id, err := parse.NewResourceID(plan.Name.ValueString(), plan.ParentID.ValueString(), plan.Type.ValueString())
	if err != nil {
		diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	client := r.ProviderData.ResourceClient
	isNewResource := responseState == nil || responseState.Raw.IsNull()
	if isNewResource {
		// check if the resource already exists
		_, err = client.Get(ctx, id.AzureResourceId, id.ApiVersion)
		if err == nil {
			diagnostics.AddError("Resource already exists", tf.ImportAsExistsError("azapi_resource", id.ID()).Error())
			return
		}
		if !utils.ResponseErrorWasNotFound(err) {
			diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("checking for presence of existing %s: %+v", id, err).Error())
			return
		}
	}

	// build the request body
	body := make(map[string]interface{})
	err = json.Unmarshal([]byte(plan.Body.ValueString()), &body)
	if err != nil {
		diagnostics.AddError("Invalid JSON string", fmt.Sprintf(`The argument "body" is invalid: value: %s, err: %+v`, plan.Body.ValueString(), err))
		return
	}
	if diagnostics.Append(expandBody(body, *plan)...); diagnostics.HasError() {
		return
	}

	if !isNewResource {
		// handle the case that identity block was once set, now it's removed
		if stateIdentity := identity.FromList(state.Identity); body["identity"] == nil && stateIdentity.Type.ValueString() != string(identity.None) {
			noneIdentity := identity.Model{Type: types.StringValue(string(identity.None))}
			out, _ := identity.ExpandIdentity(noneIdentity)
			body["identity"] = out
		}

		// handle the case that `ignore_body_changes` is set
		if ignoreChanges := AsStringList(plan.IgnoreBodyChanges); len(ignoreChanges) != 0 {
			// retrieve the existing resource
			existing, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
			if err != nil {
				diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
				return
			}

			merged, err := overrideWithPaths(body, existing, ignoreChanges)
			if err != nil {
				diagnostics.AddError("Invalid configuration", fmt.Sprintf(`The argument "ignore_body_changes" is invalid: value: %s, err: %+v`, plan.IgnoreBodyChanges.String(), err))
				return
			}

			if id.ResourceDef != nil {
				merged = (*id.ResourceDef).GetWriteOnly(utils.NormalizeObject(merged))
			}

			body = merged.(map[string]interface{})
		}
	}

	// create/update the resource
	for _, lockId := range AsStringList(plan.Locks) {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	responseBody, err := client.CreateOrUpdate(ctx, id.AzureResourceId, id.ApiVersion, body)
	if err != nil {
		diagnostics.AddError("Failed to create/update resource", fmt.Errorf("creating/updating %s: %+v", id, err).Error())
		return
	}

	// generate the computed fields
	plan.ID = types.StringValue(id.ID())
	plan.Output = types.StringValue(flattenOutput(responseBody, AsStringList(plan.ResponseExportValues)))
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		if !plan.Identity.IsNull() {
			planIdentity := identity.FromList(plan.Identity)
			if v := identity.FlattenIdentity(bodyMap["identity"]); v != nil {
				planIdentity.TenantID = v.TenantID
				planIdentity.PrincipalID = v.PrincipalID
			} else {
				planIdentity.TenantID = types.StringNull()
				planIdentity.PrincipalID = types.StringNull()
			}
			plan.Identity = identity.ToList(planIdentity)
		}
	}

	diagnostics.Append(responseState.Set(ctx, plan)...)
}

func (r *AzapiResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var model AzapiResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	id, err := parse.ResourceIDWithResourceType(model.ID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	client := r.ProviderData.ResourceClient
	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("Error reading %q - removing from state", id.ID()))
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
		return
	}

	state := model
	state.Name = types.StringValue(id.Name)
	state.ParentID = types.StringValue(id.ParentId)
	state.Type = types.StringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

	bodyJson := model.Body.ValueString()
	requestBody := make(map[string]interface{})
	err = json.Unmarshal([]byte(bodyJson), &requestBody)
	if err != nil && !model.Body.IsNull() {
		response.Diagnostics.AddError("Invalid JSON string", fmt.Sprintf(`The argument "body" is invalid: value: %s, err: %+v`, bodyJson, err))
		return
	}
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		if v, ok := bodyMap["location"]; ok && location.Normalize(v.(string)) != location.Normalize(model.Location.ValueString()) {
			state.Location = types.StringValue(v.(string))
		}
		if output := tags.FlattenTags(bodyMap["tags"]); len(output.Elements()) != 0 || len(state.Tags.Elements()) != 0 {
			state.Tags = output
		}
		if requestBody["identity"] == nil {
			identityFromResponse := identity.FlattenIdentity(bodyMap["identity"])
			switch {
			case state.Identity.IsNull() && (identityFromResponse == nil || identityFromResponse.Type.ValueString() == string(identity.None)):
				state.Identity = basetypes.NewListNull(identity.Model{}.ModelType())
			case state.Identity.IsNull() && identityFromResponse != nil && identityFromResponse.Type.ValueString() != string(identity.None):
				state.Identity = identity.ToList(*identityFromResponse)
			case !state.Identity.IsNull() && identityFromResponse == nil:
				stateIdentity := identity.FromList(state.Identity)
				if stateIdentity.Type.ValueString() == string(identity.None) {
					// do nothing
				} else {
					state.Identity = basetypes.NewListNull(identity.Model{}.ModelType())
				}
			case !state.Identity.IsNull() && identityFromResponse != nil:
				stateIdentity := identity.FromList(state.Identity)
				if len(stateIdentity.IdentityIDs.Elements()) == 0 && len(identityFromResponse.IdentityIDs.Elements()) == 0 {
					// to suppress the diff of identity_ids = [] and identity_ids = null
					identityFromResponse.IdentityIDs = stateIdentity.IdentityIDs
				}
				state.Identity = identity.ToList(*identityFromResponse)
			}
		}
	}
	state.Output = types.StringValue(flattenOutput(responseBody, AsStringList(model.ResponseExportValues)))

	if ignoreBodyChanges := AsStringList(model.IgnoreBodyChanges); len(ignoreBodyChanges) != 0 {
		if out, err := overrideWithPaths(responseBody, requestBody, ignoreBodyChanges); err == nil {
			responseBody = out
		} else {
			response.Diagnostics.AddError("Invalid configuration", fmt.Sprintf(`The argument "ignore_body_changes" is invalid: value: %s, err: %+v`, model.IgnoreBodyChanges.String(), err))
			return
		}
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
	state.Body = types.StringValue(string(data))

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (r *AzapiResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	client := r.ProviderData.ResourceClient

	var model *AzapiResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	id, err := parse.ResourceIDWithResourceType(model.ID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	for _, lockId := range AsStringList(model.Locks) {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	_, err = client.Delete(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil && !utils.ResponseErrorWasNotFound(err) {
		response.Diagnostics.AddError("Failed to delete resource", fmt.Errorf("deleting %s: %+v", id, err).Error())
	}
}

func (r *AzapiResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Importing Resource - parsing %q", request.ID))

	input := request.ID
	idUrl, err := url.Parse(input)
	if err != nil {
		response.Diagnostics.AddError("Invalid Resource ID", fmt.Errorf("parsing Resource ID %q: %+v", input, err).Error())
		return
	}
	apiVersion := idUrl.Query().Get("api-version")
	if apiVersion == "" {
		resourceType := utils.GetResourceType(input)
		apiVersions := azure.GetApiVersions(resourceType)
		if len(apiVersions) != 0 {
			input = fmt.Sprintf("%s?api-version=%s", input, apiVersions[len(apiVersions)-1])
		}
	}

	id, err := parse.ResourceIDWithApiVersion(input)
	if err != nil {
		response.Diagnostics.AddError("Invalid Resource ID", fmt.Errorf("parsing Resource ID %q: %+v", input, err).Error())
		return
	}

	client := r.ProviderData.ResourceClient

	state := AzapiResourceModel{
		ID:                      types.StringValue(id.ID()),
		Name:                    types.StringValue(id.Name),
		ParentID:                types.StringValue(id.ParentId),
		Type:                    types.StringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion)),
		Locks:                   types.ListNull(types.StringType),
		Identity:                types.ListNull(identity.Model{}.ModelType()),
		RemovingSpecialChars:    types.BoolValue(false),
		SchemaValidationEnabled: types.BoolValue(true),
		IgnoreBodyChanges:       types.ListNull(types.StringType),
		IgnoreCasing:            types.BoolValue(false),
		IgnoreMissingProperty:   types.BoolValue(true),
		ResponseExportValues:    types.ListNull(types.StringType),
		Output:                  types.StringValue("{}"),
		Tags:                    types.MapNull(types.StringType),
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("[INFO] Error reading %q - removing from state", id.ID()))
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("resource %q is imported", id.ID()))
	if id.ResourceDef != nil {
		writeOnlyBody := (*id.ResourceDef).GetWriteOnly(utils.NormalizeObject(responseBody))
		if bodyMap, ok := writeOnlyBody.(map[string]interface{}); ok {
			delete(bodyMap, "location")
			delete(bodyMap, "tags")
			delete(bodyMap, "name")
			delete(bodyMap, "identity")
			writeOnlyBody = bodyMap
		}
		data, err := json.Marshal(writeOnlyBody)
		if err != nil {
			response.Diagnostics.AddError("Invalid body", err.Error())
			return
		}
		state.Body = types.StringValue(string(data))
	} else {
		data, err := json.Marshal(responseBody)
		if err != nil {
			response.Diagnostics.AddError("Invalid body", err.Error())
			return
		}
		state.Body = types.StringValue(string(data))
	}
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		if v, ok := bodyMap["location"]; ok {
			state.Location = types.StringValue(location.Normalize(v.(string)))
		}
		if output := tags.FlattenTags(bodyMap["tags"]); len(output.Elements()) != 0 {
			state.Tags = output
		}
		if v := identity.FlattenIdentity(bodyMap["identity"]); v != nil {
			state.Identity = identity.ToList(*v)
		}
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (r *AzapiResource) nameWithDefaultNaming(name types.String, removingSpecialChars bool) (types.String, diag.Diagnostics) {
	if name.IsNull() && r.ProviderData.Features.DefaultNaming != "" {
		return types.StringValue(r.ProviderData.Features.DefaultNaming), diag.Diagnostics{}
	}

	assignedName := ""
	if !name.IsNull() && !name.IsUnknown() {
		assignedName = name.ValueString()
		if len(r.ProviderData.Features.DefaultNamingPrefix) != 0 {
			assignedName = r.ProviderData.Features.DefaultNamingPrefix + assignedName
		}
		if len(r.ProviderData.Features.DefaultNamingSuffix) != 0 {
			assignedName += r.ProviderData.Features.DefaultNamingSuffix
		}
		if removingSpecialChars {
			assignedName = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(assignedName, "")
		}
	}
	if assignedName == "" {
		return types.StringNull(), diag.Diagnostics{
			diag.NewErrorDiagnostic("Missing required argument", `The argument "name" is required, but no definition was found.`),
		}
	}
	return types.StringValue(assignedName), diag.Diagnostics{}
}

func (r *AzapiResource) tagsWithDefaultTags(config types.Map, body map[string]interface{}, state *AzapiResourceModel, resourceDef *aztypes.ResourceType) types.Map {
	if config.IsNull() {
		switch {
		case body["tags"] != nil:
			return tags.FlattenTags(body["tags"])
		case len(r.ProviderData.Features.DefaultTags) != 0 && canResourceHaveProperty(resourceDef, "tags"):
			defaultTags := r.ProviderData.Features.DefaultTags
			if state == nil || state.Tags.IsNull() {
				return tags.FlattenTags(defaultTags)
			} else {
				currentTags := tags.ExpandTags(state.Tags)
				if !reflect.DeepEqual(currentTags, defaultTags) {
					return tags.FlattenTags(defaultTags)
				} else {
					return state.Tags
				}
			}
		}
	}
	return config
}

func (r *AzapiResource) locationWithDefaultLocation(config types.String, body map[string]interface{}, state *AzapiResourceModel, resourceDef *aztypes.ResourceType) types.String {
	if config.IsNull() {
		switch {
		case body["location"] != nil:
			return types.StringValue(body["location"].(string))
		case len(r.ProviderData.Features.DefaultLocation) != 0 && canResourceHaveProperty(resourceDef, "location"):
			defaultLocation := r.ProviderData.Features.DefaultLocation
			if state == nil || state.Location.IsNull() {
				return types.StringValue(defaultLocation)
			} else {
				currentLocation := state.Location.ValueString()
				if location.Normalize(currentLocation) != location.Normalize(defaultLocation) {
					return types.StringValue(defaultLocation)
				} else {
					return state.Location
				}
			}
		}
	}
	return config
}

func expandBody(body map[string]interface{}, model AzapiResourceModel) diag.Diagnostics {
	if body == nil {
		return diag.Diagnostics{}
	}
	if body["location"] == nil && !model.Location.IsNull() && !model.Location.IsUnknown() {
		body["location"] = model.Location.ValueString()
	}
	if body["tags"] == nil && !model.Tags.IsNull() && !model.Tags.IsUnknown() {
		body["tags"] = tags.ExpandTags(model.Tags)
	}
	if body["identity"] == nil && !model.Identity.IsNull() && !model.Identity.IsUnknown() {
		identityModel := identity.FromList(model.Identity)
		out, err := identity.ExpandIdentity(identityModel)
		if err != nil {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic("Invalid configuration", fmt.Sprintf(`The argument "identity" is invalid: value: %s, err: %+v`, model.Identity.String(), err)),
			}
		}
		body["identity"] = out
	}
	return diag.Diagnostics{}
}

func validateDuplicatedDefinitions(model *AzapiResourceModel, body map[string]interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	if !model.Tags.IsNull() && !model.Tags.IsUnknown() && body["tags"] != nil {
		diags.AddError("Invalid configuration", `can't specify both the argument "tags" and "tags" in the argument "body"`)
	}
	if !model.Location.IsNull() && !model.Location.IsUnknown() && body["location"] != nil {
		diags.AddError("Invalid configuration", `can't specify both the argument "location" and "location" in the argument "body"`)
	}
	if !model.Identity.IsNull() && !model.Identity.IsUnknown() && body["identity"] != nil {
		diags.AddError("Invalid configuration", `can't specify both the argument "identity" and "identity" in the argument "body"`)
	}
	return diags
}
