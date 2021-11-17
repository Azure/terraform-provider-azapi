package utils

import (
	"net/url"
	"strings"
)

func GetId(resource interface{}) string {
	if resource == nil {
		return ""
	}
	if resourceMap, ok := resource.(map[string]interface{}); ok {
		if id, ok := resourceMap["id"]; ok {
			return id.(string)
		}
	}
	return ""
}

func GetResourceType(id string) string {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	resourceType := ""
	provider := ""
	for current := 0; current <= len(components)-2; current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return ""
		}

		if key == "providers" {
			provider = value
			resourceType = provider
		} else if len(provider) > 0 {
			resourceType += "/" + key
		}
	}
	return resourceType
}
