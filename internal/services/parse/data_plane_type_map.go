package parse

import (
	_ "embed"
	"encoding/json"
	"strings"
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
	apiPath := findApiPathByResourceType(resourceType)
	if apiPath == nil {
		return false
	}
	return strings.Contains(apiPath.UrlFormat, "{"+identifier+"}") ||
		strings.Contains(apiPath.UrlFormat, "{"+identifier+"=")
}

func HasNameSegment(resourceType string) bool {
	return hasIdentifierSegment(resourceType, "name")
}
