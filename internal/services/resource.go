package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func schemaValidation(azureResourceType, apiVersion string, resourceDef *aztypes.ResourceType, body interface{}) error {
	log.Printf("[INFO] prepare validation for resource type: %s, api-version: %s", azureResourceType, apiVersion)
	versions := azure.GetApiVersions(azureResourceType)
	if len(versions) == 0 {
		return schemaValidationError(fmt.Sprintf("the argument \"type\" is invalid.\n resource type %s can't be found.\n", azureResourceType))
	}
	isVersionValid := false
	for _, version := range versions {
		if version == apiVersion {
			isVersionValid = true
			break
		}
	}
	if !isVersionValid {
		return schemaValidationError(fmt.Sprintf("the argument \"type\"'s api-version is invalid.\n The supported versions are [%s].\n", strings.Join(versions, ", ")))
	}

	if resourceDef != nil {
		errors := (*resourceDef).Validate(utils.NormalizeObject(body), "")
		if len(errors) != 0 {
			errorMsg := "the argument \"body\" is invalid:\n"
			for _, err := range errors {
				errorMsg += fmt.Sprintf("%s\n", err.Error())
			}
			return schemaValidationError(errorMsg)
		}
	}
	return nil
}

func schemaValidationError(detail string) error {
	return fmt.Errorf("embedded schema validation failed: %s You can try to update `azapi` provider to "+
		"the latest version or disable the validation using the feature flag `schema_validation_enabled = false` "+
		"within the resource block", detail)
}

func canResourceHaveProperty(resourceDef *aztypes.ResourceType, property string) bool {
	if resourceDef == nil || resourceDef.Body == nil || resourceDef.Body.Type == nil {
		return false
	}
	objectType, ok := (*resourceDef.Body.Type).(*aztypes.ObjectType)
	if !ok {
		return false
	}
	if prop, ok := objectType.Properties[property]; ok {
		if !prop.IsReadOnly() {
			return true
		}
	}
	return false
}

func overrideWithPaths(base interface{}, changed interface{}, paths []string) (interface{}, error) {
	if len(paths) == 0 {
		return base, nil
	}

	pathSet := make(map[string]bool)
	for _, path := range paths {
		pathSet[path] = true
	}

	return utils.OverrideWithPaths(base, changed, "", pathSet)
}

func flattenOutput(responseBody interface{}, paths []string) string {
	for _, path := range paths {
		if path == "*" {
			if v, ok := responseBody.(string); ok {
				return v
			}
			outputJson, _ := json.Marshal(responseBody)
			return string(outputJson)
		}
	}

	var output interface{}
	output = make(map[string]interface{})
	for _, path := range paths {
		part := utils.ExtractObject(responseBody, path)
		if part == nil {
			continue
		}
		output = utils.MergeObject(output, part)
	}
	outputJson, _ := json.Marshal(output)
	return string(outputJson)
}

func flattenOutputPayload(responseBody interface{}, paths []string) attr.Value {
	for _, path := range paths {
		if path == "*" {
			if v, ok := responseBody.(string); ok {
				return basetypes.NewStringValue(v)
			}
			data, err := json.Marshal(responseBody)
			if err != nil {
				return nil
			}
			out, err := dynamic.FromJSONImplied(data)
			if err != nil {
				return nil
			}
			return out
		}
	}

	var output interface{}
	output = make(map[string]interface{})
	for _, path := range paths {
		part := utils.ExtractObject(responseBody, path)
		if part == nil {
			continue
		}
		output = utils.MergeObject(output, part)
	}
	data, err := json.Marshal(output)
	if err != nil {
		return nil
	}
	out, err := dynamic.FromJSONImplied(data)
	if err != nil {
		return nil
	}
	return out
}

func AsStringList(input types.List) []string {
	var result []string
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert list to string list: %s", diags))
	}
	return result
}

// bodySemanticallyEqual compares two dynamic values semantically.
// If the values are JSON strings, it normalizes the JSON strings before comparing them.
// If the values are HCL objects, it compares them semantically.
func bodySemanticallyEqual(a, b types.Dynamic) bool {
	if a.IsNull() && b.IsNull() {
		return true
	}
	if a.IsNull() || b.IsNull() {
		return false
	}
	aType := a.UnderlyingValue().Type(context.TODO())
	bType := b.UnderlyingValue().Type(context.TODO())
	if aType.Equal(types.StringType) && bType.Equal(types.StringType) {
		aJson := a.UnderlyingValue().(types.String).ValueString()
		bJson := b.UnderlyingValue().(types.String).ValueString()
		return utils.NormalizeJson(aJson) == utils.NormalizeJson(bJson)
	}
	return dynamic.SemanticallyEqual(a, b)
}

func unmarshalBody(input types.Dynamic, out interface{}) error {
	if input.IsNull() || input.IsUnknown() || input.IsUnderlyingValueUnknown() {
		return nil
	}
	switch value := input.UnderlyingValue().(type) {
	case types.String:
		err := json.Unmarshal([]byte(value.ValueString()), out)
		if err != nil {
			return fmt.Errorf(`unmarshaling failed: value: %s, err: %+v`, value.ValueString(), err)
		}
		return nil
	default:
		data, err := dynamic.ToJSON(input)
		if err != nil {
			return fmt.Errorf(`invalid dynamic value: value: %s, err: %+v`, input.String(), err)
		}
		if err = json.Unmarshal(data, &out); err != nil {
			return fmt.Errorf(`unmarshaling failed: value: %s, err: %+v`, string(data), err)
		}
		return nil
	}
}

func dynamicIsString(input types.Dynamic) bool {
	if input.IsNull() || input.IsUnknown() || input.IsUnderlyingValueUnknown() {
		return false
	}
	_, ok := input.UnderlyingValue().(types.String)
	return ok
}
