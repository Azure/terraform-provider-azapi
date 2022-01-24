package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
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

	resourceDef, err := azure.GetResourceDefinition(id.AzureResourceType, id.ApiVersion)
	if err == nil && resourceDef != nil {
		errors := (*resourceDef).Validate(utils.NormalizeObject(body), "")
		if len(errors) != 0 {
			errorMsg := "the `body` is invalid: \n"
			for _, err := range errors {
				errorMsg += fmt.Sprintf("%s\n", err.Error())
			}
			return fmt.Errorf(errorMsg)
		}
	} else {
		log.Printf("[ERROR] load embedded schema: %+v\n", err)
	}
	return nil
}
