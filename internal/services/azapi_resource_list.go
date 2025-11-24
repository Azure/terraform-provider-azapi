package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AzapiResourceListModel struct {
	Type            types.String `tfsdk:"type"`
	ParentID        types.String `tfsdk:"parent_id"`
	Headers         types.Map    `tfsdk:"headers"`
	QueryParameters types.Map    `tfsdk:"query_parameters"`
}

// resourceListConfigValidator validates the configuration for azapi_resource_list
type resourceListConfigValidator struct{}

func (v resourceListConfigValidator) Description(_ context.Context) string {
	return "Validates that when 'type' is omitted, 'parent_id' must be a resource group ID"
}

func (v resourceListConfigValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v resourceListConfigValidator) ValidateListResourceConfig(ctx context.Context, req list.ValidateConfigRequest, resp *list.ValidateConfigResponse) {
	var config AzapiResourceListModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If type is null or empty, validate that parent_id is a resource group
	if config.Type.ValueString() == "" {
		parentIdResourceType := utils.GetResourceType(config.ParentID.ValueString())
		if !strings.EqualFold(parentIdResourceType, arm.ResourceGroupResourceType.String()) {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				fmt.Sprintf("When 'type' is omitted, 'parent_id' must be a resource group ID. Got resource type: %s (expected: %s)",
					parentIdResourceType, arm.ResourceGroupResourceType.String()),
			)
		}
	}
}

type AzapiResourceList struct {
	ProviderData *clients.Client
}

var _ list.ListResource = &AzapiResourceList{}
var _ list.ListResourceWithConfigure = &AzapiResourceList{}
var _ list.ListResourceWithConfigValidators = &AzapiResourceList{}

func (r *AzapiResourceList) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource"
}

func (r *AzapiResourceList) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	} else {
		response.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *clients.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
}

func (r *AzapiResourceList) ListResourceConfigValidators(_ context.Context) []list.ConfigValidator {
	return []list.ConfigValidator{
		&resourceListConfigValidator{},
	}
}

func (r *AzapiResourceList) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Configuration for listing Azure Resource Manager resources",
		MarkdownDescription: "This list resource allows you to list Azure Resource Manager resources of a specific type under a given scope.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type() + " When omitted, `parent_id` must be a resource group ID and all resources in the resource group will be listed.",
			},

			"parent_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: docstrings.ParentID(),
			},

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
	}
}

func (r *AzapiResourceList) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	var model AzapiResourceListModel

	diags := request.Config.Get(ctx, &model)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var listUrl string
	var apiVersion string

	// Check if type is provided
	if model.Type.IsNull() || model.Type.ValueString() == "" {
		// Type is omitted - list all resources in the resource group
		// Note: Validation that parent_id is a resource group is handled by the config validator
		listUrl = fmt.Sprintf("%s/resources", strings.TrimSuffix(model.ParentID.ValueString(), "/"))
		apiVersion = "2025-04-01"

		ctx = tflog.SetField(ctx, "parent_id", model.ParentID.ValueString())
		ctx = tflog.SetField(ctx, "list_all_resources", true)
	} else {
		// Type is provided - use existing logic
		id, err := parse.NewResourceIDSkipScopeValidation("", model.ParentID.ValueString(), model.Type.ValueString())
		if err != nil {
			diags.AddError("Invalid configuration", err.Error())
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}

		listUrl = strings.TrimSuffix(id.AzureResourceId, "/")
		apiVersion = id.ApiVersion

		ctx = tflog.SetField(ctx, "resource_type", model.Type.ValueString())
		ctx = tflog.SetField(ctx, "parent_id", model.ParentID.ValueString())
	}

	// Prepare request options
	client := r.ProviderData.ResourceClient
	requestOptions := clients.RequestOptions{
		Headers:         common.AsMapOfString(model.Headers),
		QueryParameters: clients.NewQueryParameters(common.AsMapOfLists(model.QueryParameters)),
	}

	// Make the list request
	responseBody, err := client.List(ctx, listUrl, apiVersion, requestOptions)
	if err != nil {
		diags.AddError("Failed to list resources", fmt.Sprintf("Failed to list resources, url: %s, error: %s", listUrl, err.Error()))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	responseMap, ok := responseBody.(map[string]interface{})
	if !ok {
		diags.AddError("Invalid response format", "Response is not a valid JSON object")
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if responseMap["value"] == nil {
		diags.AddError("No resources found", "The response does not contain any resources")
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	valueArray, ok := responseMap["value"].([]interface{})
	if !ok {
		diags.AddError("Invalid response format", "'value' field is not an array")
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	// Create an iterator that yields ListResult for each resource
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range valueArray {
			result := request.NewListResult(ctx)

			itemMap, ok := item.(map[string]interface{})
			if !ok {
				result.Diagnostics.AddError("Invalid item format", "Resource item is not a valid object")
				if !push(result) {
					return
				}
				continue
			}

			// Extract resource ID for identity
			resourceID, ok := itemMap["id"].(string)
			if !ok || resourceID == "" {
				result.Diagnostics.AddError("Missing resource ID", "Resource item does not contain a valid 'id' field")
				continue
			}

			id, err := parse.ResourceID(resourceID)
			if err != nil {
				result.Diagnostics.AddError("Invalid resource ID", fmt.Sprintf("Resource ID %q is invalid: %v", resourceID, err))
				continue
			}

			itemType := model.Type.ValueString()
			if itemType == "" {
				itemType = fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion)
			}

			result.DisplayName = fmt.Sprintf("%s - %s", itemType, id.Name)

			// Set the identity using the model
			identityData := AzapiResourceIdentityModel{
				ID:   types.StringValue(resourceID),
				Type: types.StringValue(itemType),
			}
			result.Diagnostics.Append(result.Identity.Set(ctx, identityData)...)

			// If full resource data is requested, populate it
			if request.IncludeResource {
				state := NewDefaultAzapiResourceModel()
				state.ID = types.StringValue(id.ID())
				state.Name = types.StringValue(id.Name)
				state.ParentID = types.StringValue(id.ParentId)
				state.Type = types.StringValue(fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

				// Use the itemMap from the list response directly instead of making another GET call
				if model.Type.ValueString() != "" {
					responseBody := itemMap
					tflog.Info(ctx, fmt.Sprintf("resource %q is populated from list response", id.ID()))
					payload, err := flattenBody(responseBody, id.ResourceDef)
					if err != nil {
						result.Diagnostics.AddError("Invalid body", err.Error())
					} else {
						state.Body = payload

						if v, ok := responseBody["location"]; ok && v != nil {
							state.Location = types.StringValue(location.Normalize(v.(string)))
						}
						if output := tags.FlattenTags(responseBody["tags"]); len(output.Elements()) != 0 {
							state.Tags = output
						}
						if v := identity.FlattenIdentity(responseBody["identity"]); v != nil {
							state.Identity = identity.ToList(*v)
						}
					}
				} else {
					// When listing all resources in a resource group, the list API doesn't return full details
					// Make a GET request to fetch the complete resource data
					tflog.Debug(ctx, fmt.Sprintf("Fetching full details for resource %q via GET request", id.ID()))

					responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.DefaultRequestOptions())
					if err != nil {
						tflog.Warn(ctx, fmt.Sprintf("Failed to fetch full details for resource %q: %v", id.ID(), err))
						result.Diagnostics.AddWarning(
							"Failed to fetch resource details",
							fmt.Sprintf("Could not retrieve full details for resource %q: %v", id.ID(), err),
						)
					} else {
						if responseMap, ok := responseBody.(map[string]interface{}); ok {
							payload, err := flattenBody(responseMap, id.ResourceDef)
							if err != nil {
								result.Diagnostics.AddError("Invalid body", err.Error())
							} else {
								state.Body = payload

								if v, ok := responseMap["location"]; ok && v != nil {
									state.Location = types.StringValue(location.Normalize(v.(string)))
								}
								if output := tags.FlattenTags(responseMap["tags"]); len(output.Elements()) != 0 {
									state.Tags = output
								}
								if v := identity.FlattenIdentity(responseMap["identity"]); v != nil {
									state.Identity = identity.ToList(*v)
								}
							}
						} else {
							result.Diagnostics.AddWarning("Invalid response format", "GET response is not a valid JSON object")
						}
					}
				}

				result.Diagnostics.Append(result.Resource.Set(ctx, state)...)
			}

			// Push the result to the stream
			if !push(result) {
				return
			}
		}
	}
}
