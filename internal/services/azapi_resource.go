package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/defaults"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/migration"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier/planmodifierdynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
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

type AzapiResourceModel struct {
	Body                          types.Dynamic       `tfsdk:"body"`
	ID                            types.String        `tfsdk:"id"`
	Identity                      types.List          `tfsdk:"identity"`
	IgnoreCasing                  types.Bool          `tfsdk:"ignore_casing"`
	IgnoreMissingProperty         types.Bool          `tfsdk:"ignore_missing_property"`
	Location                      types.String        `tfsdk:"location"`
	Locks                         types.List          `tfsdk:"locks"`
	Name                          types.String        `tfsdk:"name"`
	Output                        types.Dynamic       `tfsdk:"output"`
	ParentID                      types.String        `tfsdk:"parent_id"`
	ReplaceTriggersExternalValues types.Dynamic       `tfsdk:"replace_triggers_external_values"`
	ReplaceTriggersRefs           types.List          `tfsdk:"replace_triggers_refs"`
	ResponseExportValues          types.Dynamic       `tfsdk:"response_export_values"`
	Retry                         retry.RetryValue    `tfsdk:"retry"`
	SchemaValidationEnabled       types.Bool          `tfsdk:"schema_validation_enabled"`
	Tags                          types.Map           `tfsdk:"tags"`
	Timeouts                      timeouts.Value      `tfsdk:"timeouts"`
	Type                          types.String        `tfsdk:"type"`
	CreateHeaders                 map[string]string   `tfsdk:"create_headers"`
	CreateQueryParameters         map[string][]string `tfsdk:"create_query_parameters"`
	UpdateHeaders                 map[string]string   `tfsdk:"update_headers"`
	UpdateQueryParameters         map[string][]string `tfsdk:"update_query_parameters"`
	DeleteHeaders                 map[string]string   `tfsdk:"delete_headers"`
	DeleteQueryParameters         map[string][]string `tfsdk:"delete_query_parameters"`
	ReadHeaders                   map[string]string   `tfsdk:"read_headers"`
	ReadQueryParameters           map[string][]string `tfsdk:"read_query_parameters"`
}

var _ resource.Resource = &AzapiResource{}
var _ resource.ResourceWithConfigure = &AzapiResource{}
var _ resource.ResourceWithModifyPlan = &AzapiResource{}
var _ resource.ResourceWithValidateConfig = &AzapiResource{}
var _ resource.ResourceWithImportState = &AzapiResource{}
var _ resource.ResourceWithUpgradeState = &AzapiResource{}

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

func (r *AzapiResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: migration.AzapiResourceMigrationV0ToV2(ctx),
		1: migration.AzapiResourceMigrationV1ToV2(ctx),
	}
}

func (r *AzapiResource) Schema(ctx context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This resource can manage any Azure Resource Manager resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: docstrings.Type(),
			},

			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Specifies the name of the azure resource. Changing this forces a new resource to be created.",
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
				MarkdownDescription: docstrings.ParentID(),
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type(),
			},

			"location": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					myplanmodifier.UseStateWhen(func(a, b types.String) bool {
						return location.Normalize(a.ValueString()) == location.Normalize(b.ValueString())
					}),
				},
				MarkdownDescription: docstrings.Location(),
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
			},

			"replace_triggers_external_values": schema.DynamicAttribute{
				Optional: true,
				MarkdownDescription: "Will trigger a replace of the resource when the value changes and is not `null`. This can be used by practitioners to force a replace of the resource when certain values change, e.g. changing the SKU of a virtual machine based on the value of variables or locals. " +
					"The value is a `dynamic`, so practitioners can compose the input however they wish. For a \"break glass\" set the value to `null` to prevent the plan modifier taking effect. \n" +
					"If you have `null` values that you do want to be tracked as affecting the resource replacement, include these inside an object. \n" +
					"Advanced use cases are possible and resource replacement can be triggered by values external to the resource, for example when a dependent resource changes.\n\n" +
					"e.g. to replace a resource when either the SKU or os_type attributes change:\n" +
					"\n" +
					"```hcl\n" +
					"resource \"azapi_resource\" \"example\" {\n" +
					"  name      = var.name\n" +
					"  type      = \"Microsoft.Network/publicIPAddresses@2023-11-01\"\n" +
					"  parent_id = \"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example\"\n" +
					"  body      = {\n" +
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

			"ignore_casing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.BoolDefault(false),
				MarkdownDescription: docstrings.IgnoreCasing(),
			},

			"ignore_missing_property": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.BoolDefault(true),
				MarkdownDescription: docstrings.IgnoreMissingProperty(),
			},

			"response_export_values": CommonAttributeResponseExportValues(),

			"locks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
				MarkdownDescription: docstrings.Locks(),
			},

			"schema_validation_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.BoolDefault(true),
				MarkdownDescription: docstrings.SchemaValidationEnabled(),
			},

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("azapi_resource"),
			},

			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.Map{
					tags.Validator(),
				},
				MarkdownDescription: "A mapping of tags which should be assigned to the Azure resource.",
			},

			"retry": retry.SingleNestedAttribute(ctx),

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
			"identity": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Validators: []validator.Object{myvalidator.IdentityValidator()},
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{stringvalidator.OneOf(
								string(identity.SystemAssignedUserAssigned),
								string(identity.UserAssigned),
								string(identity.SystemAssigned),
								string(identity.None),
							)},
							MarkdownDescription: docstrings.IdentityType(),
						},

						"identity_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(myvalidator.StringIsUserAssignedIdentityID()),
							},
							MarkdownDescription: docstrings.IdentityIds(),
						},

						"principal_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: docstrings.IdentityPrincipalID(),
						},

						"tenant_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: docstrings.IdentityTenantID(),
						},
					},
				},
			},

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

	if !dynamic.IsFullyKnown(config.Body) {
		return
	}

	body := make(map[string]interface{})
	if err := unmarshalBody(config.Body, &body); err != nil {
		response.Diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
		return
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

	// Output is a computed field, it defaults to unknown if there's any plan change
	// It sets to the state if the state exists, and will set to unknown if the output needs to be updated
	if state != nil {
		plan.Output = state.Output
	}
	resourceType := config.Type.ValueString()

	// for resource group, if parent_id is not specified, set it to subscription id
	if config.ParentID.IsNull() {
		azureResourceType, _, _ := utils.GetAzureResourceTypeApiVersion(resourceType)
		if strings.EqualFold(azureResourceType, arm.ResourceGroupResourceType.String()) {
			plan.ParentID = types.StringValue(fmt.Sprintf("/subscriptions/%s", r.ProviderData.Account.GetSubscriptionId()))
		}
	}

	if name, diags := r.nameWithDefaultNaming(config.Name); !diags.HasError() {
		plan.Name = name
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

	if !dynamic.IsFullyKnown(plan.Body) {
		if config.Tags.IsNull() {
			plan.Tags = basetypes.NewMapUnknown(types.StringType)
		}
		if config.Location.IsNull() {
			plan.Location = basetypes.NewStringUnknown()
		}
		plan.Output = basetypes.NewDynamicUnknown()
		return
	}

	if state == nil || !plan.Identity.Equal(state.Identity) || !plan.ResponseExportValues.Equal(state.ResponseExportValues) || !dynamic.SemanticallyEqual(plan.Body, state.Body) {
		plan.Output = basetypes.NewDynamicUnknown()
	}

	body := make(map[string]interface{})
	if err := unmarshalBody(config.Body, &body); err != nil {
		response.Diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
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
		plan.Output = basetypes.NewDynamicUnknown()
	}

	// location field has a field level plan modifier which suppresses the diff if the location is not actually changed
	locationValue := plan.Location
	// For the following cases, we need to use the location in config as the specified location
	// case 1. To create a new resource, the location is not specified in config, then the planned location will be unknown
	// case 2. To update a resource, the location is not specified in config, then the planned location will be the state location
	if locationValue.IsUnknown() || config.Location.IsNull() {
		locationValue = config.Location
	}
	// locationWithDefaultLocation will return the location in config if it's not null, otherwise it will return the default location if it supports location
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

	// Check if any paths in replace_triggers_refs have changed
	if state != nil && plan != nil && !plan.ReplaceTriggersRefs.IsNull() {
		refPaths := make(map[string]string)
		for pathIndex, refPath := range AsStringList(plan.ReplaceTriggersRefs) {
			refPaths[fmt.Sprintf("%d", pathIndex)] = refPath
		}

		// read previous values from state
		var data interface{}
		err = json.Unmarshal([]byte(state.Body.String()), &data)
		if err != nil {
			response.Diagnostics.AddError("Invalid state body configuration", err.Error())
			return
		}
		previousValues := flattenOutputJMES(data, refPaths)

		// read current values from plan
		err = json.Unmarshal([]byte(plan.Body.String()), &data)
		if err != nil {
			response.Diagnostics.AddError("Invalid plan body configuration", err.Error())
			return
		}
		currentValues := flattenOutputJMES(data, refPaths)

		// compare previous and current values
		if !reflect.DeepEqual(previousValues, currentValues) {
			response.RequiresReplace.Append(path.Root("body"))
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

	var client clients.Requester
	client = r.ProviderData.ResourceClient
	if !plan.Retry.IsNull() {
		bkof, regexps := clients.NewRetryableErrors(
			plan.Retry.GetIntervalSeconds(),
			plan.Retry.GetMaxIntervalSeconds(),
			plan.Retry.GetMultiplier(),
			plan.Retry.GetRandomizationFactor(),
			plan.Retry.GetErrorMessageRegex(),
		)
		client = r.ProviderData.ResourceClient.WithRetry(bkof, regexps)
	}
	isNewResource := responseState == nil || responseState.Raw.IsNull()

	var timeout time.Duration
	var diags diag.Diagnostics
	if isNewResource {
		timeout, diags = plan.Timeouts.Create(ctx, 30*time.Minute)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return
		}
	} else {
		timeout, diags = plan.Timeouts.Update(ctx, 30*time.Minute)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return
		}
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if isNewResource {
		// check if the resource already exists using the non-retry client to avoid issue where user specifies
		// a FooResourceNotFound error as a retryable error
		_, err = r.ProviderData.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(plan.ReadHeaders, plan.ReadQueryParameters))
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
	if err := unmarshalBody(plan.Body, &body); err != nil {
		diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
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
	}

	// create/update the resource
	for _, lockId := range AsStringList(plan.Locks) {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	options := clients.NewRequestOptions(plan.CreateHeaders, plan.CreateQueryParameters)
	if !isNewResource {
		options = clients.NewRequestOptions(plan.UpdateHeaders, plan.UpdateQueryParameters)
	}
	_, err = client.CreateOrUpdate(ctx, id.AzureResourceId, id.ApiVersion, body, options)
	if err != nil {
		if isNewResource {
			if responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(plan.ReadHeaders, plan.ReadQueryParameters)); err == nil {
				// generate the computed fields
				plan.ID = types.StringValue(id.ID())

				output, err := buildOutputFromBody(responseBody, plan.ResponseExportValues)
				if err != nil {
					diagnostics.AddError("Failed to build output", err.Error())
					return
				}
				plan.Output = output

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
		}
		diagnostics.AddError("Failed to create/update resource", fmt.Errorf("creating/updating %s: %+v", id, err).Error())
		return
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(plan.ReadHeaders, plan.ReadQueryParameters))
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("Error reading %q - removing from state", id.ID()))
			responseState.RemoveResource(ctx)
			return
		}
		diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("reading %s: %+v", id, err).Error())
		return
	}

	// generate the computed fields
	plan.ID = types.StringValue(id.ID())

	output, err := buildOutputFromBody(responseBody, plan.ResponseExportValues)
	if err != nil {
		diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	plan.Output = output

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

	readTimeout, diags := model.Timeouts.Read(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	id, err := parse.ResourceIDWithResourceType(model.ID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
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
		client = r.ProviderData.ResourceClient.WithRetry(bkof, regexps)
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(model.ReadHeaders, model.ReadQueryParameters))
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

	requestBody := make(map[string]interface{})
	if err := unmarshalBody(model.Body, &requestBody); err != nil {
		response.Diagnostics.AddError("Invalid body", fmt.Sprintf(`The argument "body" is invalid: %s`, err.Error()))
		return
	}

	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		if v, ok := bodyMap["location"]; ok && v != nil && location.Normalize(v.(string)) != location.Normalize(model.Location.ValueString()) {
			state.Location = types.StringValue(v.(string))
		}
		if output := tags.FlattenTags(bodyMap["tags"]); len(output.Elements()) != 0 || len(state.Tags.Elements()) != 0 {
			state.Tags = output
		}
		if requestBody["identity"] == nil {
			// The following codes are used to reflect the actual changes of identity when it's not configured inside the body.
			// And it suppresses the diff of nil identity and identity whose type is none.
			identityFromResponse := identity.FlattenIdentity(bodyMap["identity"])
			switch {
			// Identity is not specified in config, and it's not in the response
			case state.Identity.IsNull() && (identityFromResponse == nil || identityFromResponse.Type.ValueString() == string(identity.None)):
				state.Identity = basetypes.NewListNull(identity.Model{}.ModelType())

			// Identity is not specified in config, but it's in the response
			case state.Identity.IsNull() && identityFromResponse != nil && identityFromResponse.Type.ValueString() != string(identity.None):
				state.Identity = identity.ToList(*identityFromResponse)

			// Identity is specified in config, but it's not in the response
			case !state.Identity.IsNull() && identityFromResponse == nil:
				stateIdentity := identity.FromList(state.Identity)
				// skip when the configured identity type is none
				if stateIdentity.Type.ValueString() == string(identity.None) {
					// do nothing
				} else {
					state.Identity = basetypes.NewListNull(identity.Model{}.ModelType())
				}

			// Identity is specified in config, and it's in the response
			case !state.Identity.IsNull() && identityFromResponse != nil:
				stateIdentity := identity.FromList(state.Identity)
				// suppress the diff of identity_ids = [] and identity_ids = null
				if len(stateIdentity.IdentityIDs.Elements()) == 0 && len(identityFromResponse.IdentityIDs.Elements()) == 0 {
					// to suppress the diff of identity_ids = [] and identity_ids = null
					identityFromResponse.IdentityIDs = stateIdentity.IdentityIDs
				}
				state.Identity = identity.ToList(*identityFromResponse)
			}
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

	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	state.Output = output

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
		state.Body = payload
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (r *AzapiResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var model *AzapiResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
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
		client = r.ProviderData.ResourceClient.WithRetry(bkof, regexps)
	}

	deleteTimeout, diags := model.Timeouts.Delete(ctx, 30*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	id, err := parse.ResourceIDWithResourceType(model.ID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	for _, lockId := range AsStringList(model.Locks) {
		locks.ByID(lockId)
		defer locks.UnlockByID(lockId)
	}

	_, err = client.Delete(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(model.DeleteHeaders, model.DeleteQueryParameters))
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
		Body:                    types.DynamicNull(),
		SchemaValidationEnabled: types.BoolValue(true),
		IgnoreCasing:            types.BoolValue(false),
		IgnoreMissingProperty:   types.BoolValue(true),
		ResponseExportValues:    types.DynamicNull(),
		Output:                  types.DynamicNull(),
		Tags:                    types.MapNull(types.StringType),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"read":   types.StringType,
				"delete": types.StringType,
			}),
		},
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.NewRequestOptions(state.ReadHeaders, state.ReadQueryParameters))
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
		payload, err := dynamic.FromJSONImplied(data)
		if err != nil {
			response.Diagnostics.AddError("Invalid payload", err.Error())
			return
		}
		state.Body = payload
	} else {
		data, err := json.Marshal(responseBody)
		if err != nil {
			response.Diagnostics.AddError("Invalid body", err.Error())
			return
		}
		payload, err := dynamic.FromJSONImplied(data)
		if err != nil {
			response.Diagnostics.AddError("Invalid payload", err.Error())
			return
		}
		state.Body = payload
	}
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		if v, ok := bodyMap["location"]; ok && v != nil {
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

func (r *AzapiResource) nameWithDefaultNaming(config types.String) (types.String, diag.Diagnostics) {
	if !config.IsNull() {
		return config, diag.Diagnostics{}
	}
	if r.ProviderData.Features.DefaultNaming != "" {
		return types.StringValue(r.ProviderData.Features.DefaultNaming), diag.Diagnostics{}
	}
	return types.StringNull(), diag.Diagnostics{
		diag.NewErrorDiagnostic("Missing required argument", `The argument "name" is required, but no definition was found.`),
	}
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
		// To suppress the diff of config: tags = null and state: tags = {}
		case state != nil && !state.Tags.IsUnknown() && len(state.Tags.Elements()) == 0:
			return state.Tags
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
		// To suppress the diff of config: location = null and state: location = ""
		// This case happens when upgrading resources which doesn't support location from terraform-plugin-sdk built azapi provider
		case state != nil && !state.Location.IsUnknown() && state.Location.ValueString() == "":
			return state.Location
		}
	}
	return config
}

func expandBody(body map[string]interface{}, model AzapiResourceModel) diag.Diagnostics {
	if body == nil {
		return diag.Diagnostics{}
	}
	if body["location"] == nil && !model.Location.IsNull() && !model.Location.IsUnknown() && len(model.Location.ValueString()) != 0 {
		body["location"] = model.Location.ValueString()
	}
	if body["tags"] == nil && !model.Tags.IsNull() && !model.Tags.IsUnknown() && len(model.Tags.Elements()) != 0 {
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
