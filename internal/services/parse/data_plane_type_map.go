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
