package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jmespath/go-jmespath"
)

const (
	readActionNone     = "none"
	readActionListName = "list"
	readActionPOSTVerb = "POST"
)

// resolveReadResponse returns the response body that should be used for the read-merge. When a
// plain GET covers none of the configured body's leaves, some Azure endpoints only expose their
// values through a POST list function; readAction controls whether that fallback is used.
func resolveReadResponse(ctx context.Context, client *clients.ResourceClient, azureResourceId, azureResourceType, apiVersion string, requestBody map[string]interface{}, getResponse interface{}, readAction, readActionMethod, readActionResponsePath string, requestOptions clients.RequestOptions) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	action := strings.TrimSpace(readAction)

	if action == "" {
		return autoReadResponse(ctx, client, azureResourceId, azureResourceType, apiVersion, requestBody, getResponse, readActionResponsePath, requestOptions)
	}

	if strings.EqualFold(action, readActionNone) {
		return getResponse, diags
	}

	method := strings.TrimSpace(readActionMethod)
	if method == "" {
		method = readActionPOSTVerb
	}

	response, err := client.Action(ctx, azureResourceId, action, apiVersion, method, nil, requestOptions)
	if err != nil {
		diags.AddError("Failed to read resource via action", fmt.Sprintf("%s %q on %q failed: %s", method, action, azureResourceId, err.Error()))
		return getResponse, diags
	}
	return applyResponsePath(response, readActionResponsePath, &diags), diags
}

func autoReadResponse(ctx context.Context, client *clients.ResourceClient, azureResourceId, azureResourceType, apiVersion string, requestBody map[string]interface{}, getResponse interface{}, readActionResponsePath string, requestOptions clients.RequestOptions) (interface{}, diag.Diagnostics) {
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

	response = applyResponsePath(response, readActionResponsePath, &diags)
	if diags.HasError() {
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
		diags.AddError("Invalid read_action_response_path", fmt.Sprintf("failed to evaluate JMESPath %q against the read action response: %s", path, err.Error()))
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
