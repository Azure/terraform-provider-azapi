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

func flattenOutput(responseBody interface{}, paths []string) attr.Value {
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

func flattenOutputJMES(responseBody interface{}, paths map[string]string) attr.Value {
	var output interface{}
	output = make(map[string]interface{})
	for pathKey, path := range paths {
		part := utils.ExtractObjectJMES(responseBody, pathKey, path)
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

func unmarshalBody(input types.Dynamic, out interface{}) error {
	if input.IsNull() || input.IsUnknown() || input.IsUnderlyingValueUnknown() {
		return nil
	}
	data, err := dynamic.ToJSON(input)
	if err != nil {
		return fmt.Errorf(`invalid dynamic value: value: %s, err: %+v`, input.String(), err)
	}
	if err = json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf(`unmarshaling failed: value: %s, err: %+v`, string(data), err)
	}
	return nil
}
