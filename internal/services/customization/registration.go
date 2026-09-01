package customization

import "strings"

var customizations = make(map[string]DataPlaneResource)

func init() {
	var keyVaultKeyCustomization DataPlaneResource = KeyVaultKeyCustomization{}
	customizations[strings.ToLower(keyVaultKeyCustomization.GetResourceType())] = keyVaultKeyCustomization

	var keyVaultCertificateCustomization DataPlaneResource = KeyVaultCertificateCustomization{}
	customizations[strings.ToLower(keyVaultCertificateCustomization.GetResourceType())] = keyVaultCertificateCustomization

	var foundryAgentCustomization DataPlaneResource = FoundryAgentCustomization{}
	customizations[strings.ToLower(foundryAgentCustomization.GetResourceType())] = foundryAgentCustomization

	var foundryEvaluationCustomization DataPlaneResource = FoundryEvaluationCustomization{}
	customizations[strings.ToLower(foundryEvaluationCustomization.GetResourceType())] = foundryEvaluationCustomization

	var foundryEvaluationRunCustomization DataPlaneResource = FoundryEvaluationRunCustomization{}
	customizations[strings.ToLower(foundryEvaluationRunCustomization.GetResourceType())] = foundryEvaluationRunCustomization
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
