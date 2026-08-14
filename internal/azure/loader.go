package azure

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

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

// skipApiVersions contains resource type + API version combinations that are present
// in the schema but not yet supported by Azure, causing "NoRegisteredProviderFound" errors at runtime.
// Key: resource type (case-insensitive, stored lowercase), Value: set of API versions to skip.
var skipApiVersions = map[string]map[string]bool{
	"microsoft.keyvault/vaults/keys":    {"2026-02-01": true, "2026-03-01-preview": true},
	"microsoft.keyvault/vaults/secrets": {"2026-02-01": true, "2026-03-01-preview": true},
}

func GetApiVersions(resourceType string) []string {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return []string{}
	}
	skipped := skipApiVersions[strings.ToLower(resourceType)]
	res := make([]string, 0)
	for key, value := range azureSchema.Resources {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				if skipped[v.ApiVersion] {
					continue
				}
				res = append(res, v.ApiVersion)
			}
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

// GetFunctionDefinitions returns all resource function definitions registered for the given
// resource type and api-version (e.g. list, listKeys). Returns an empty slice when none exist.
func GetFunctionDefinitions(resourceType, apiVersion string) ([]*types.ResourceFunctionType, error) {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return nil, fmt.Errorf("failed to load azure schema index")
	}
	res := make([]*types.ResourceFunctionType, 0)
	for key, value := range azureSchema.Functions {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				if v.ApiVersion != apiVersion {
					continue
				}
				def, err := v.GetDefinition()
				if err != nil {
					return nil, err
				}
				if def != nil {
					res = append(res, def)
				}
			}
		}
	}
	return res, nil
}

// GetFunctionDefinition returns the resource function definition matching the given name
// (case-insensitive) for the resource type and api-version, or nil when not found.
func GetFunctionDefinition(resourceType, apiVersion, name string) (*types.ResourceFunctionType, error) {
	definitions, err := GetFunctionDefinitions(resourceType, apiVersion)
	if err != nil {
		return nil, err
	}
	for _, def := range definitions {
		if strings.EqualFold(def.Name, name) {
			return def, nil
		}
	}
	return nil, nil
}
