package parse

import (
	_ "embed"
	"encoding/json"
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

func hasIdentifierSegment(resourceType string, identifier string) bool {
	if azureResourceType, _, err := utils.GetAzureResourceTypeApiVersion(resourceType); err == nil {
		resourceType = azureResourceType
	}
	apiPath := findApiPathByResourceType(resourceType)
	if apiPath == nil {
		// Unknown type: assume it needs an identifier so validation isn't silently skipped.
		return true
	}
	return strings.Contains(apiPath.UrlFormat, "{"+identifier+"}") ||
		strings.Contains(apiPath.UrlFormat, "{"+identifier+"=")
}

func HasNameSegment(resourceType string) bool {
	return hasIdentifierSegment(resourceType, "name")
}
