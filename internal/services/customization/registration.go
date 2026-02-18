package customization

import "strings"

var customizations = make(map[string]DataPlaneResource)

func init() {
	var keyVaultKeyCustomization DataPlaneResource = KeyVaultKeyCustomization{}
	customizations[strings.ToLower(keyVaultKeyCustomization.GetResourceType())] = keyVaultKeyCustomization

	var aiFoundryAssistantCustomization DataPlaneResource = AIFoundryAssistantCustomization{}
	customizations[strings.ToLower(aiFoundryAssistantCustomization.GetResourceType())] = aiFoundryAssistantCustomization
}

func GetCustomization(resourceType string) *DataPlaneResource {
	// remove api-version
	resourceType = strings.Split(resourceType, "@")[0]
	resourceType = strings.ToLower(resourceType)
	customization, exists := customizations[resourceType]
	if !exists {
		return nil
	}
	return &customization
}
