package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jmespath/go-jmespath"
)

const (
	readActionListName = "list"
	readActionPOSTVerb = "POST"
)

type ReadOverrideModel struct {
	Method       types.String `tfsdk:"method"`
	Action       types.String `tfsdk:"action"`
	ResponsePath types.String `tfsdk:"response_path"`
}

func readOverrideAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"method":        types.StringType,
		"action":        types.StringType,
		"response_path": types.StringType,
	}
}

func readOverrideSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"method": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The HTTP method used to read the resource. The only supported value is `POST`.",
				Validators: []validator.String{
					stringvalidator.OneOf(readActionPOSTVerb),
				},
			},
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the action appended to the resource ID, for example `list`.",
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
			},
			"response_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional [JMESPath](https://jmespath.org/) expression applied to the action response before it is merged into `body`.",
			},
		},
		MarkdownDescription: "Overrides the request used to read resource values when a plain `GET` does not return them. When unset, the provider automatically attempts a `POST {resource_id}/list` fallback when supported and the `GET` response contains none of the configured `body` values.",
	}
}

func readOverrideFromObject(ctx context.Context, value types.Object) (*ReadOverrideModel, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var model ReadOverrideModel
	diags := value.As(ctx, &model, basetypes.ObjectAsOptions{})
	return &model, diags
}

func resolveReadResponse(ctx context.Context, client *clients.ResourceClient, azureResourceId, azureResourceType, apiVersion string, requestBody map[string]interface{}, getResponse interface{}, readOverride *ReadOverrideModel, requestOptions clients.RequestOptions) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if readOverride == nil {
		return autoReadResponse(ctx, client, azureResourceId, azureResourceType, apiVersion, requestBody, getResponse, requestOptions)
	}

	action := strings.TrimSpace(readOverride.Action.ValueString())
	method := strings.TrimSpace(readOverride.Method.ValueString())

	response, err := client.Action(ctx, azureResourceId, action, apiVersion, method, nil, requestOptions)
	if err != nil {
		diags.AddError("Failed to read resource via action", fmt.Sprintf("%s %q on %q failed: %s", method, action, azureResourceId, err.Error()))
		return getResponse, diags
	}
	return applyResponsePath(response, readOverride.ResponsePath.ValueString(), &diags), diags
}

func autoReadResponse(ctx context.Context, client *clients.ResourceClient, azureResourceId, azureResourceType, apiVersion string, requestBody map[string]interface{}, getResponse interface{}, requestOptions clients.RequestOptions) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	total, covered := countCoveredLeaves(requestBody, getResponse)
	if total == 0 || covered > 0 {
		return getResponse, diags
	}

	fn, err := azure.GetFunctionDefinition(azureResourceType, apiVersion, readActionListName)
	if err != nil || fn == nil {
		return getResponse, diags
	}

	response, err := client.Action(ctx, azureResourceId, fn.Name, apiVersion, readActionPOSTVerb, nil, requestOptions)
	if err != nil {
		diags.AddWarning("Failed to read resource via list action", fmt.Sprintf("The GET response for %q did not contain any configured values, so a POST %q was attempted but failed: %s. Falling back to the GET response.", azureResourceId, fn.Name, err.Error()))
		return getResponse, diags
	}

	if _, postCovered := countCoveredLeaves(requestBody, response); postCovered == 0 {
		diags.AddWarning("List action did not cover configured values", fmt.Sprintf("The list function %q for %q did not return any of the configured values. Falling back to the GET response.", fn.Name, azureResourceId))
		return getResponse, diags
	}
	return response, diags
}

func applyResponsePath(response interface{}, path string, diags *diag.Diagnostics) interface{} {
	path = strings.TrimSpace(path)
	if path == "" {
		return response
	}
	result, err := jmespath.Search(path, response)
	if err != nil {
		diags.AddError("Invalid read_override response_path", fmt.Sprintf("failed to evaluate JMESPath %q against the read override response: %s", path, err.Error()))
		return response
	}
	return result
}

// countCoveredLeaves mirrors the read-merge: it walks the config body and counts how many leaf
// values have a corresponding non-nil value at the same path in response. covered == 0 means the
// response would contribute nothing to the merge.
func countCoveredLeaves(body interface{}, response interface{}) (total int, covered int) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		if response != nil {
			return 1, 1
		}
		return 1, 0
	}
	responseMap, _ := response.(map[string]interface{})
	for key, value := range bodyMap {
		if childMap, isMap := value.(map[string]interface{}); isMap && len(childMap) > 0 {
			var childResponse interface{}
			if responseMap != nil {
				childResponse = responseMap[key]
			}
			t, c := countCoveredLeaves(childMap, childResponse)
			total += t
			covered += c
			continue
		}
		total++
		if responseMap != nil && responseMap[key] != nil {
			covered++
		}
	}
	return total, covered
}
