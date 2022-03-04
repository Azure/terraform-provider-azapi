package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
)

func schemaValidation(id parse.ResourceId, body interface{}) error {
	log.Printf("[INFO] prepare validation for resource type: %s, api-version: %s", id.AzureResourceType, id.ApiVersion)
	versions := azure.GetApiVersions(id.AzureResourceType)
	if len(versions) == 0 {
		return fmt.Errorf("the `type` is invalid, resource type %s can't be found", id.AzureResourceType)
	}
	isVersionValid := false
	for _, version := range versions {
		if version == id.ApiVersion {
			isVersionValid = true
			break
		}
	}
	if !isVersionValid {
		return fmt.Errorf("the `type`'s api-version is invalid. The supported versions are [%s]\n", strings.Join(versions, ", "))
	}

	if id.ResourceDef != nil {
		errors := (*id.ResourceDef).Validate(utils.NormalizeObject(body), "")
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
