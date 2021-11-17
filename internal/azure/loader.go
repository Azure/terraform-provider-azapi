package azure

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/types"
)

var schema *Schema

//go:embed generated
var StaticFiles embed.FS

func GetAzureSchema() *Schema {
	if schema == nil {
		data, err := StaticFiles.ReadFile("generated/index.json")
		if err != nil {
			log.Printf("[ERROR] failed to load schema index: %+v", err)
			return nil
		}
		err = json.Unmarshal(data, &schema)
		if err != nil {
			log.Printf("[ERROR] failed to unmarshal schema index: %+v", err)
			return nil
		}
	}
	return schema
}

func GetApiVersions(resourceType string) []string {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return []string{}
	}
	res := make([]string, 0)
	if azureSchema.Resources[resourceType] != nil {
		for _, v := range azureSchema.Resources[resourceType].Definitions {
			res = append(res, v.ApiVersion)
		}
	}
	sort.Strings(res)
	return res
}

func GetResourceDefinition(resourceType, apiVersion string) (*types.ResourceType, error) {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return nil, fmt.Errorf("failed to load azure schema index")
	}
	if azureSchema.Resources[resourceType] != nil {
		for _, v := range azureSchema.Resources[resourceType].Definitions {
			if v.ApiVersion == apiVersion {
				return v.GetDefinition()
			}
		}
	}
	return nil, fmt.Errorf("failed to find resource type %s api-version %s in azure schema index", resourceType, apiVersion)
}
