package parse

import (
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/utils"
)

type DataPlaneResourceId struct {
	AzureResourceId   string
	ApiVersion        string
	AzureResourceType string
	Name              string
	ParentId          string
}

func NewDataPlaneResourceId(name, parentId, resourceType string) (DataPlaneResourceId, error) {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return DataPlaneResourceId{}, err
	}

	azureResourceId := ""
	if apiPath := findApiPathByResourceType(azureResourceType); apiPath != nil {
		parts := strings.Split(apiPath.UrlFormat, "/")
		for i, part := range parts {
			switch {
			case part == "{parentId}":
				parts[i] = parentId
			case part == "{name}":
				parts[i] = name
			case part == "{apiVersion}":
				parts[i] = apiVersion
			case strings.HasPrefix(part, "{name="):
				defaultName := part[6 : len(part)-1]
				if !strings.EqualFold(name, defaultName) {
					return DataPlaneResourceId{}, fmt.Errorf("name %s is not equal to %s", name, defaultName)
				}
				parts[i] = defaultName
			case strings.Contains(part, "{name}"):
				// Handle embedded {name} placeholder, e.g., "indexes('{name}')"
				parts[i] = strings.ReplaceAll(part, "{name}", name)
			}
		}
		azureResourceId = strings.Join(parts, "/")
	}

	return DataPlaneResourceId{
		AzureResourceId:   azureResourceId,
		ApiVersion:        apiVersion,
		AzureResourceType: azureResourceType,
		Name:              name,
		ParentId:          parentId,
	}, nil
}

// DataPlaneResourceIDWithResourceType parses a Resource ID and resource type into an ResourceId struct
func DataPlaneResourceIDWithResourceType(azureResourceId, resourceType string) (DataPlaneResourceId, error) {
	azureResourceType, _, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return DataPlaneResourceId{}, err
	}

	name := ""
	parentId := ""
	if apiPath := findApiPathByResourceType(azureResourceType); apiPath != nil {
		urlFormatParts := strings.Split(apiPath.UrlFormat, "/")
		azureResourceIdParts := strings.Split(azureResourceId, "/")
		j := len(azureResourceIdParts) - 1
		for i := len(urlFormatParts) - 1; i >= 0; i-- {
			if j < 0 {
				return DataPlaneResourceId{}, fmt.Errorf("index %d is less than 0", j)
			}
			switch {
			case strings.HasPrefix(urlFormatParts[i], "{name"):
				name = azureResourceIdParts[j]
				j--
			case strings.Contains(urlFormatParts[i], "{name}"):
				// Handle embedded {name} placeholder, e.g., "indexes('{name}')"
				// Extract the actual name from patterns like "indexes('myindex')"
				template := urlFormatParts[i]
				actual := azureResourceIdParts[j]
				// Find the position of {name} in the template
				nameStart := strings.Index(template, "{name}")
				if nameStart >= 0 {
					// Extract prefix and suffix around {name}
					prefix := template[:nameStart]
					suffix := template[nameStart+6:] // 6 is len("{name}")
					// Remove prefix and suffix from actual value to get the name
					if strings.HasPrefix(actual, prefix) && strings.HasSuffix(actual, suffix) {
						name = actual[len(prefix) : len(actual)-len(suffix)]
					}
				}
				j--
			case urlFormatParts[i] == "{parentId}":
				for j >= 0 {
					if j > 0 && i > 0 && azureResourceIdParts[j-1] == urlFormatParts[i-1] {
						break
					}
					parentId = azureResourceIdParts[j] + "/" + parentId
					j--
				}
			case urlFormatParts[i] == "{apiVersion}":
				j--
			case strings.EqualFold(azureResourceIdParts[j], urlFormatParts[i]):
				j--
			}
		}
	}
	parentId = strings.TrimSuffix(parentId, "/")

	return NewDataPlaneResourceId(name, parentId, resourceType)
}

func (id DataPlaneResourceId) String() string {
	segments := []string{
		fmt.Sprintf("ResourceId %q", id.AzureResourceId),
		fmt.Sprintf("Api Version %q", id.ApiVersion),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource", segmentsStr)
}

func (id DataPlaneResourceId) ID() string {
	return id.AzureResourceId
}
