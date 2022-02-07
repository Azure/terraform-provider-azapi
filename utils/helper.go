package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/types"
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
	if len(id) == 0 || id == "/" {
		return "Tenant"
	}
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}
	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	if len(components) == 2 && strings.EqualFold(components[0], "subscriptions") {
		return "Subscription"
	}
	if len(components) == 4 &&
		strings.EqualFold(components[0], "providers") &&
		strings.EqualFold(components[1], "Microsoft.Management") &&
		strings.EqualFold(components[2], "managementGroups") {
		return "Microsoft.Management/managementGroups"
	}
	if len(components) == 4 &&
		strings.EqualFold(components[0], "subscriptions") &&
		strings.EqualFold(components[2], "resourceGroups") {
		return "Microsoft.Resources/resourceGroups"
	}
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

func GetName(id string) string {
	if index := strings.LastIndex(id, "/"); index != -1 {
		return id[index+1:]
	}
	return ""
}

func GetParentId(id string) string {
	if len(id) == 0 || id == "/" {
		return ""
	}
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}
	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	if len(components) == 2 && strings.EqualFold(components[0], "subscriptions") {
		return ""
	}
	if len(components) == 4 &&
		strings.EqualFold(components[0], "providers") &&
		strings.EqualFold(components[1], "Microsoft.Management") &&
		strings.EqualFold(components[2], "managementGroups") {
		return ""
	}
	if len(components) == 4 &&
		strings.EqualFold(components[0], "subscriptions") &&
		strings.EqualFold(components[2], "resourceGroups") {
		return fmt.Sprintf("/%s/%s", components[0], components[1])
	}
	if firstIndex := strings.Index(path, "providers"); firstIndex > 0 {
		if lastIndex := strings.LastIndex(path, "providers"); lastIndex != -1 {
			if firstIndex != lastIndex {
				return "/" + path[0:lastIndex-1]
			}
		}
	}
	parentId := ""
	for current := 0; current <= len(components)-4; current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return ""
		}

		if current == len(components)-4 && key == "providers" {

		} else {
			parentId += "/" + key + "/" + value
		}
	}
	return parentId
}

func GetParentType(resourceType string) string {
	parts := strings.Split(resourceType, "/")
	if len(parts) <= 2 {
		return ""
	}
	return strings.Join(parts[0:len(parts)-1], "/")
}

func GetScopeType(id string) types.ScopeType {
	if len(id) == 0 || id == "/" {
		return types.Tenant
	}
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return types.Unknown
	}
	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	componentsBeforeProvider := make([]string, 0)
	for _, part := range components {
		if strings.EqualFold(part, "providers") {
			break
		}
		componentsBeforeProvider = append(componentsBeforeProvider, part)
	}

	if len(componentsBeforeProvider) == 0 {
		mgmtPrefix := "providers/Microsoft.Management/managementGroups"
		if len(path) >= len(mgmtPrefix) && strings.EqualFold(path[0:len(mgmtPrefix)], mgmtPrefix) {
			return types.ManagementGroup
		} else {
			return types.Tenant
		}
	}

	if len(componentsBeforeProvider) == 2 && strings.EqualFold(componentsBeforeProvider[0], "subscriptions") {
		return types.Subscription
	}

	if len(componentsBeforeProvider) == 4 && strings.EqualFold(componentsBeforeProvider[0], "subscriptions") &&
		strings.EqualFold(componentsBeforeProvider[2], "resourceGroups") {
		return types.ResourceGroup
	}

	return types.Unknown
}

func ExpandStringSlice(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return &result
}
