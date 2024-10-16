package azure

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
)

var schema *Schema

//go:embed generated
var StaticFiles embed.FS

var mutex = &sync.Mutex{}

func GetAzureSchema() *Schema {
	mutex.Lock()
	defer mutex.Unlock()
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
	for key, value := range azureSchema.Resources {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				res = append(res, v.ApiVersion)
			}
		}
	}

	// TODO: remove the below codes when Resources RP 2024-07-01 is available
	if strings.EqualFold(resourceType, arm.ResourceGroupResourceType.String()) {
		temp := make([]string, 0)
		for _, v := range res {
			if v != "2024-07-01" {
				temp = append(temp, v)
			}
		}
		res = temp
	}

	sort.Strings(res)
	return res
}

func GetResourceDefinition(resourceType, apiVersion string) (*types.ResourceType, error) {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return nil, fmt.Errorf("failed to load azure schema index")
	}
	for key, value := range azureSchema.Resources {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				if v.ApiVersion == apiVersion {
					return v.GetDefinition()
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find resource type %s api-version %s in azure schema index", resourceType, apiVersion)
}
