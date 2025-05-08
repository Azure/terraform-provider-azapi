package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/utils"
)

var template string

type ResourceProvider struct {
	ResourceProviderNamespace string `json:"resourceProviderNamespace"`
	FriendlyName              string `json:"friendlyName"`
}

type ResourceType struct {
	ResourceType string `json:"resourceType"`
	FriendlyName string `json:"friendlyName"`
}

var (
	resourceProviders map[string]ResourceProvider
	resourceTypes     map[string]ResourceType
)

func init() {
	var items []ResourceProvider
	mappingJsonPath := path.Join("tools", "generator-example-doc", "resource_providers.json")
	// #nosec G304
	data, err := os.ReadFile(mappingJsonPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &items)
	if err != nil {
		panic(err)
	}

	resourceProviders = make(map[string]ResourceProvider)
	for _, item := range items {
		resourceProviders[strings.ToLower(item.ResourceProviderNamespace)] = item
	}

	var resourceTypeItems []ResourceType
	resourceTypeJsonPath := path.Join("tools", "generator-example-doc", "resource_types.json")
	// #nosec G304
	data, err = os.ReadFile(resourceTypeJsonPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &resourceTypeItems)
	if err != nil {
		panic(err)
	}
	resourceTypes = make(map[string]ResourceType)
	for _, item := range resourceTypeItems {
		resourceTypes[strings.ToLower(item.ResourceType)] = item
	}

	data, err = os.ReadFile(path.Join("tools", "generator-example-doc", "template.md"))
	if err != nil {
		panic(err)
	}
	template = string(data)
}

func main() {
	inputDir := flag.String("input-dir", "./examples", "directory to scan for example files")
	outputDir := flag.String("output-dir", "./docs/guides", "directory to write documentation files")

	flag.Parse()
	if *inputDir == "" || *outputDir == "" {
		log.Fatal("input-dir and output-dir flags are required")
	}

	resourceTypeDirs, err := os.ReadDir(*inputDir)
	if err != nil {
		log.Fatalf("Error reading input directory: %s", err)
	}

	for _, resourceTypeDir := range resourceTypeDirs {
		if !resourceTypeDir.IsDir() {
			continue
		}
		if !strings.Contains(resourceTypeDir.Name(), "@") {
			continue
		}

		content, err := generateDocumentation(path.Join(*inputDir, resourceTypeDir.Name()))
		if err != nil {
			log.Fatalf("Error generating documentation for %s: %s", resourceTypeDir.Name(), err)
		}

		resourceTypeName := strings.Split(resourceTypeDir.Name(), "@")[0]
		outputFile := path.Join(*outputDir, resourceTypeName+".md")
		// #nosec G306
		err = os.WriteFile(outputFile, []byte(content), 0644)
		if err != nil {
			log.Printf("Error writing documentation for %s: %s", resourceTypeDir.Name(), err)
		}
	}
}

func generateDocumentation(inputDir string) (string, error) {
	resourceType := strings.Split(path.Base(inputDir), "@")[0]
	resourceType = strings.ReplaceAll(resourceType, "_", "/")

	resourceProviderName := strings.Split(resourceType, "/")[0]
	resourceTypeWithoutRP := strings.Join(strings.Split(resourceType, "/")[1:], "/")

	resourceProviderFriendlyName := GetResourceProviderFriendlyName(resourceProviderName)
	if resourceProviderFriendlyName == "" {
		return "", fmt.Errorf("resource provider %s friendly name not found, please add it to resource_providers.json", resourceProviderName)
	}

	resourceTypeFriendlyName := GetResourceTypeFriendlyName(resourceType)
	if resourceTypeFriendlyName == "" {
		return "", fmt.Errorf("resource type %s friendly name not found, please add it to resource_types.json", resourceType)
	}

	apiVersions := azure.GetApiVersions(resourceType)
	apiVersion := "API_VERSION"
	if len(apiVersions) > 0 {
		apiVersion = apiVersions[len(apiVersions)-1]
	}

	parentIds := getParentIds(resourceType)

	resourceId := ""
	if len(parentIds) > 0 {
		lastSegment := resourceType[strings.LastIndex(resourceType, "/")+1:]
		if utils.IsTopLevelResourceType(resourceType) {
			resourceId = fmt.Sprintf("%s/providers/%s/%s/{resourceName}", parentIds[0], resourceProviderName, lastSegment)
		} else {
			resourceId = fmt.Sprintf("%s/%s/{resourceName}", parentIds[0], lastSegment)
		}
	}

	out := template
	out = strings.ReplaceAll(out, "{{.subcategory}}", fmt.Sprintf("%s - %s", resourceProviderName, resourceProviderFriendlyName))
	out = strings.ReplaceAll(out, "{{.page_title}}", resourceTypeWithoutRP)
	out = strings.ReplaceAll(out, "{{.resource_type}}", resourceType)
	out = strings.ReplaceAll(out, "{{.resource_type_friendly_name}}", resourceTypeFriendlyName)
	out = strings.ReplaceAll(out, "{{.reference_link}}", fmt.Sprintf("https://learn.microsoft.com/en-us/azure/templates/%s?pivots=deployment-language-terraform", resourceType))
	out = strings.ReplaceAll(out, "{{.api_versions}}", strings.Join(addBackticks(apiVersions), ", "))
	out = strings.ReplaceAll(out, "{{.api_version}}", apiVersion)
	out = strings.ReplaceAll(out, "{{.resource_id}}", resourceId)
	out = strings.ReplaceAll(out, "{{.parent_id}}", strings.Join(addBackticks(parentIds), "  \n  "))

	// key is the scenario name, value is the example content
	exampleMap := make(map[string]string)
	scenarioDirs, err := os.ReadDir(inputDir)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}
	for _, scenarioDir := range scenarioDirs {
		if !scenarioDir.IsDir() || scenarioDir.Name() == "testdata" {
			continue
		}

		scenarioName := scenarioDir.Name()
		exampleFilePath := path.Join(inputDir, scenarioName, "main.tf")
		// #nosec G304
		exampleContent, err := os.ReadFile(exampleFilePath)
		if err != nil {
			log.Printf("Error reading example file for %s: %s", exampleFilePath, err)
			continue
		}

		exampleMap[scenarioName] = string(exampleContent)
	}
	// check if there's main.tf in the inputDir
	mainFilePath := path.Join(inputDir, "main.tf")
	if _, err := os.Stat(mainFilePath); err == nil {
		// #nosec G304
		exampleContent, err := os.ReadFile(mainFilePath)
		if err != nil {
			log.Printf("Error reading example file for %s: %s", mainFilePath, err)
			return "", err
		}
		exampleMap["default"] = string(exampleContent)
	}

	scenarioNames := make([]string, 0)
	for scenarioName := range exampleMap {
		scenarioNames = append(scenarioNames, scenarioName)
	}
	slices.Sort(scenarioNames)

	example := ""
	for _, scenarioName := range scenarioNames {
		exampleContent := exampleMap[scenarioName]
		example += fmt.Sprintf("### %s\n\n", scenarioName)
		example += fmt.Sprintf("```hcl\n%s\n```\n\n", exampleContent)
	}
	out = strings.ReplaceAll(out, "{{.example}}", example)

	return out, nil
}

func getParentIds(resourceType string) []string {
	const defaultParentId = "{any azure resource id}"
	apiVersions := azure.GetApiVersions(resourceType)
	if len(apiVersions) == 0 {
		return nil
	}
	resourceDef, err := azure.GetResourceDefinition(resourceType, apiVersions[0])
	if err != nil || resourceDef == nil {
		return nil
	}

	if utils.IsTopLevelResourceType(resourceType) {
		scopeIds := make([]string, 0)
		for _, scope := range resourceDef.ScopeTypes {
			switch scope {
			case types.Tenant:
				scopeIds = append(scopeIds, "/")
			case types.Subscription:
				scopeIds = append(scopeIds, "/subscriptions/{subscriptionId}")
			case types.ManagementGroup:
				scopeIds = append(scopeIds, "/providers/Microsoft.Management/managementGroups/{managementGroupId}")
			case types.ResourceGroup:
				scopeIds = append(scopeIds, "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}")
			case types.Extension:
				scopeIds = append(scopeIds, defaultParentId)
			default:
				scopeIds = append(scopeIds, defaultParentId)
			}
		}
		return scopeIds
	}

	parentResourceType := resourceType[:strings.LastIndex(resourceType, "/")]
	parentResourceParentIds := getParentIds(parentResourceType)
	if len(parentResourceParentIds) == 0 {
		return []string{defaultParentId}
	}

	parentIds := make([]string, 0)
	lastSegment := parentResourceType[strings.LastIndex(parentResourceType, "/")+1:]
	resourceProvider := strings.Split(parentResourceType, "/")[0]
	for _, parentId := range parentResourceParentIds {
		if parentId == defaultParentId {
			parentIds = append(parentIds, defaultParentId)
			continue
		}
		if utils.IsTopLevelResourceType(parentResourceType) {
			parentIds = append(parentIds, fmt.Sprintf("%s/providers/%s/%s/{resourceName}", parentId, resourceProvider, lastSegment))
		} else {
			parentIds = append(parentIds, fmt.Sprintf("%s/%s/{resourceName}", parentId, lastSegment))
		}
	}

	return parentIds
}

func GetResourceProviderFriendlyName(resourceProviderName string) string {
	if resourceProvider, ok := resourceProviders[strings.ToLower(resourceProviderName)]; ok {
		return resourceProvider.FriendlyName
	}
	return ""
}

func GetResourceTypeFriendlyName(resourceType string) string {
	if v, ok := resourceTypes[strings.ToLower(resourceType)]; ok {
		return v.FriendlyName
	}
	return ""
}

func addBackticks(input []string) []string {
	for i, str := range input {
		if str != "" {
			input[i] = fmt.Sprintf("`%s`", str)
		}
	}
	return input
}
