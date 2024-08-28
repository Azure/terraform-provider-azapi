package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
	tf_types "github.com/hashicorp/terraform-plugin-framework/types"
)

func GetId(resource interface{}) *string {
	if resource == nil {
		return nil
	}
	if resourceMap, ok := resource.(map[string]interface{}); ok {
		if id, ok := resourceMap["id"]; ok {
			idStr := ""
			if id != nil {
				idStr = id.(string)
			}
			return &idStr
		}
	}
	return nil
}

func GetResourceType(id string) string {
	if id == "/" {
		return arm.TenantResourceType.String()
	}
	resourceType, err := arm.ParseResourceType(id)
	if err != nil {
		return ""
	}
	return resourceType.String()
}

func GetName(id string) string {
	resourceId, err := arm.ParseResourceID(id)
	if err != nil {
		return ""
	}
	return resourceId.Name
}

func GetParentId(id string) string {
	resourceId, err := arm.ParseResourceID(id)
	if err != nil {
		return ""
	}
	if resourceId.Parent.ResourceType.String() == arm.TenantResourceType.String() {
		return "/"
	}
	return resourceId.Parent.String()
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

	// #nosec G602
	if len(componentsBeforeProvider) == 2 && strings.EqualFold(componentsBeforeProvider[0], "subscriptions") {
		return types.Subscription
	}

	// #nosec G602
	if len(componentsBeforeProvider) == 4 && strings.EqualFold(componentsBeforeProvider[0], "subscriptions") &&
		// #nosec G602
		strings.EqualFold(componentsBeforeProvider[2], "resourceGroups") {
		return types.ResourceGroup
	}

	return types.Unknown
}

func GetAzureResourceTypeApiVersion(resourceType string) (string, string, error) {
	parts := strings.Split(resourceType, "@")
	apiVersion := ""
	azureResourceType := ""
	if len(parts) == 2 {
		apiVersion = parts[1]
		azureResourceType = parts[0]
	} else {
		return "", "", fmt.Errorf("`type` is invalid, it should be like `ResourceProvider/resourceTypes@ApiVersion`")
	}
	return azureResourceType, apiVersion, nil
}

func IsTopLevelResourceType(resourceType string) bool {
	return len(strings.Split(resourceType, "/")) == 2
}

func GetAzureResourceType(resourceType string, apiVersion string) string {
	return fmt.Sprintf("%s@%s", resourceType, apiVersion)
}

func GetAzureResourceTypeParts(azureResourceType string) (resourceProvider string, parts []string, err error) {
	parts = strings.Split(azureResourceType, "/")
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid azure resource type, it should be like `ResourceProvider/resourceTypes`")
	}
	return parts[0], parts[1:], nil
}

func ParseResourceNames(resourceNamesParam tf_types.List) []string {
	resourceNamesList := resourceNamesParam.Elements()

	resourceNames := make([]string, len(resourceNamesList))
	for i, element := range resourceNamesList {
		resourceNames[i] = element.(tf_types.String).ValueString()
	}

	return resourceNames
}

// TryAppendDefaultApiVersion appends the default api version to the resource type if it is not already present.
func TryAppendDefaultApiVersion(resourceType string) string {
	if !strings.Contains(resourceType, "@") {
		resourceType += "@latest"
	}
	return resourceType
}
