package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
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

type AzapiResourceList struct {
	ProviderData *clients.Client
}

var _ list.ListResource = &AzapiResourceList{}
var _ list.ListResourceWithConfigure = &AzapiResourceList{}

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

func (r *AzapiResourceList) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Configuration for listing Azure Resource Manager resources",
		MarkdownDescription: "This list resource allows you to list Azure Resource Manager resources of a specific type under a given scope.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type(),
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

	// Parse the resource ID
	id, err := parse.NewResourceIDSkipScopeValidation("", model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		diags.AddError("Invalid configuration", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	ctx = tflog.SetField(ctx, "resource_type", model.Type.ValueString())
	ctx = tflog.SetField(ctx, "parent_id", model.ParentID.ValueString())

	listUrl := strings.TrimSuffix(id.AzureResourceId, "/")

	// Prepare request options
	client := r.ProviderData.ResourceClient
	requestOptions := clients.RequestOptions{
		Headers:         common.AsMapOfString(model.Headers),
		QueryParameters: clients.NewQueryParameters(common.AsMapOfLists(model.QueryParameters)),
	}

	// Make the list request
	responseBody, err := client.List(ctx, listUrl, id.ApiVersion, requestOptions)
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
				result.Diagnostics.AddWarning("Missing resource ID", "Resource item does not contain a valid 'id' field")
				resourceID = ""
			}

			id, err := parse.ResourceIDWithResourceType(resourceID, model.Type.ValueString())
			if err != nil {
				result.Diagnostics.AddWarning("Invalid resource ID", fmt.Sprintf("Resource ID %q is invalid: %s", resourceID, err.Error()))
				// Set a generic display name if parsing fails
				result.DisplayName = resourceID
			} else {
				result.DisplayName = fmt.Sprintf("%s - %s", model.Type.ValueString(), id.Name)
			}

			// Set the identity using the model
			identityData := AzapiResourceIdentityModel{
				ID:   types.StringValue(resourceID),
				Type: model.Type,
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

				result.Diagnostics.Append(result.Resource.Set(ctx, state)...)
			}

			// Push the result to the stream
			if !push(result) {
				return
			}
		}
	}
}
