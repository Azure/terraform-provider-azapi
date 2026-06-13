package parse

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Azure/terraform-provider-azapi/utils"
)

type ApiPath struct {
	UrlFormat       string
	ResourceType    string
	URL             string
	ParentIDExample string
}

//go:embed data_plane_resources.json
var raw string

var apiPaths = make([]ApiPath, 0)

func init() {
	err := json.Unmarshal([]byte(raw), &apiPaths)
	if err != nil {
		panic(err)
	}
}

func findApiPathByResourceType(resourceType string) *ApiPath {
	for _, apiPath := range apiPaths {
		if strings.EqualFold(apiPath.ResourceType, resourceType) {
			return &apiPath
		}
	}
	return nil
}

var dataPlanePlaceholderPattern = regexp.MustCompile(`\{([^{}]+)\}`)

type dataPlanePlaceholder struct {
	Key          string
	DefaultValue string
	HasDefault   bool
}

func placeholdersForURLFormat(urlFormat string) []dataPlanePlaceholder {
	matches := dataPlanePlaceholderPattern.FindAllStringSubmatch(urlFormat, -1)
	result := make([]dataPlanePlaceholder, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		parts := strings.SplitN(match[1], "=", 2)
		placeholder := dataPlanePlaceholder{
			Key: strings.TrimSpace(parts[0]),
		}
		if len(parts) == 2 {
			placeholder.HasDefault = true
			placeholder.DefaultValue = strings.TrimSpace(parts[1])
		}
		result = append(result, placeholder)
	}
	return result
}

func DataPlaneResourcePlaceholderKeys(resourceType string) ([]string, error) {
	azureResourceType, _, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return nil, err
	}
	apiPath := findApiPathByResourceType(azureResourceType)
	if apiPath == nil {
		return nil, fmt.Errorf("unsupported data plane resource type %q", azureResourceType)
	}

	keys := make([]string, 0)
	for _, placeholder := range placeholdersForURLFormat(apiPath.UrlFormat) {
		if placeholder.Key == "parentId" || placeholder.Key == "apiVersion" {
			continue
		}
		if !slices.Contains(keys, placeholder.Key) {
			keys = append(keys, placeholder.Key)
		}
	}
	return keys, nil
}
