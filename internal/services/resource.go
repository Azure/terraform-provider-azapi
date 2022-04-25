package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/utils"
)

func schemaValidation(azureResourceType, apiVersion string, resourceDef *types.ResourceType, body interface{}) error {
	log.Printf("[INFO] prepare validation for resource type: %s, api-version: %s", azureResourceType, apiVersion)
	versions := azure.GetApiVersions(azureResourceType)
	if len(versions) == 0 {
		return fmt.Errorf("the `type` is invalid, resource type %s can't be found", azureResourceType)
	}
	isVersionValid := false
	for _, version := range versions {
		if version == apiVersion {
			isVersionValid = true
			break
		}
	}
	if !isVersionValid {
		return fmt.Errorf("the `type`'s api-version is invalid. The supported versions are [%s]\n", strings.Join(versions, ", "))
	}

	if resourceDef != nil {
		errors := (*resourceDef).Validate(utils.NormalizeObject(body), "")
		if len(errors) != 0 {
			errorMsg := "the `body` is invalid: \n"
			for _, err := range errors {
				errorMsg += fmt.Sprintf("%s\n", err.Error())
			}
			return fmt.Errorf(errorMsg)
		}
	}
	return nil
}

func isResourceHasProperty(resourceDef *types.ResourceType, property string) bool {
	if resourceDef == nil || resourceDef.Body == nil || resourceDef.Body.Type == nil {
		return false
	}
	objectType, ok := (*resourceDef.Body.Type).(*types.ObjectType)
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

func flattenOutput(responseBody interface{}, paths []interface{}) string {
	for _, path := range paths {
		if path == "*" {
			outputJson, _ := json.Marshal(responseBody)
			return string(outputJson)
		}
	}

	var output interface{}
	output = make(map[string]interface{})
	for _, path := range paths {
		part := utils.ExtractObject(responseBody, path.(string))
		if part == nil {
			continue
		}
		output = utils.GetMergedJson(output, part)
	}
	outputJson, _ := json.Marshal(output)
	return string(outputJson)
}
