package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/terraform-provider-azapi/utils"
)

type DataPlaneResourceId struct {
	AzureResourceId   string
	ApiVersion        string
	AzureResourceType string
	Name              string
	ParentId          string
	Identifiers       map[string]string
}

func NewDataPlaneResourceId(name, parentId, resourceType string) (DataPlaneResourceId, error) {
	return NewDataPlaneResourceIdWithIdentifiers(name, parentId, resourceType, nil)
}

func NewDataPlaneResourceIdWithIdentifiers(name, parentId, resourceType string, identifiers map[string]string) (DataPlaneResourceId, error) {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return DataPlaneResourceId{}, err
	}

	azureResourceId := ""
	if apiPath := findApiPathByResourceType(azureResourceType); apiPath != nil {
		values := map[string]string{
			"parentId":   parentId,
			"apiVersion": apiVersion,
		}
		if name != "" {
			values["name"] = name
		}
		for key, value := range identifiers {
			values[key] = value
		}

		azureResourceId, err = renderDataPlaneURLFormat(apiPath.UrlFormat, values)
		if err != nil {
			return DataPlaneResourceId{}, err
		}
	}

	return DataPlaneResourceId{
		AzureResourceId:   azureResourceId,
		ApiVersion:        apiVersion,
		AzureResourceType: azureResourceType,
		Name:              name,
		ParentId:          parentId,
		Identifiers:       cloneStringMap(identifiers),
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
	identifiers := make(map[string]string)
	if apiPath := findApiPathByResourceType(azureResourceType); apiPath != nil {
		values, err := parseDataPlaneURLFormat(apiPath.UrlFormat, azureResourceId)
		if err != nil {
			return DataPlaneResourceId{}, err
		}
		parentId = values["parentId"]
		name = values["name"]
		for key, value := range values {
			if key == "parentId" || key == "name" || key == "apiVersion" {
				continue
			}
			identifiers[key] = value
		}
	}
	parentId = strings.TrimSuffix(parentId, "/")

	return NewDataPlaneResourceIdWithIdentifiers(name, parentId, resourceType, identifiers)
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

func renderDataPlaneURLFormat(urlFormat string, values map[string]string) (string, error) {
	missing := make([]string, 0)
	rendered := dataPlanePlaceholderPattern.ReplaceAllStringFunc(urlFormat, func(token string) string {
		match := dataPlanePlaceholderPattern.FindStringSubmatch(token)
		if len(match) < 2 {
			return token
		}
		parts := strings.SplitN(match[1], "=", 2)
		key := strings.TrimSpace(parts[0])
		if len(parts) == 2 {
			defaultValue := strings.TrimSpace(parts[1])
			if value, ok := values[key]; ok && value != "" && !strings.EqualFold(value, defaultValue) {
				missing = append(missing, fmt.Sprintf("%s must equal %s", key, defaultValue))
				return token
			}
			return defaultValue
		}
		value, ok := values[key]
		if !ok || value == "" {
			missing = append(missing, key)
			return token
		}
		return value
	})
	if len(missing) != 0 {
		return "", fmt.Errorf("missing required data plane identifiers for URL format %q: %s", urlFormat, strings.Join(missing, ", "))
	}
	return rendered, nil
}

func parseDataPlaneURLFormat(urlFormat string, actual string) (map[string]string, error) {
	regexText := "^"
	lastIndex := 0
	matches := dataPlanePlaceholderPattern.FindAllStringSubmatchIndex(urlFormat, -1)
	for _, match := range matches {
		start := match[0]
		end := match[1]
		contentStart := match[2]
		contentEnd := match[3]

		regexText += regexp.QuoteMeta(urlFormat[lastIndex:start])

		content := urlFormat[contentStart:contentEnd]
		parts := strings.SplitN(content, "=", 2)
		if len(parts) == 2 {
			regexText += regexp.QuoteMeta(strings.TrimSpace(parts[1]))
		} else {
			if end == len(urlFormat) {
				regexText += "(.+)"
			} else {
				regexText += "(.+?)"
			}
		}
		lastIndex = end
	}
	regexText += regexp.QuoteMeta(urlFormat[lastIndex:]) + "$"

	re, err := regexp.Compile(regexText)
	if err != nil {
		return nil, fmt.Errorf("building parser for data plane URL format %q: %w", urlFormat, err)
	}
	submatches := re.FindStringSubmatch(actual)
	if submatches == nil {
		return nil, fmt.Errorf("resource ID %q does not match data plane URL format %q", actual, urlFormat)
	}

	result := make(map[string]string)
	groupIndex := 1
	for _, placeholder := range placeholdersForURLFormat(urlFormat) {
		if placeholder.HasDefault {
			result[placeholder.Key] = placeholder.DefaultValue
			continue
		}
		if groupIndex >= len(submatches) {
			return nil, fmt.Errorf("resource ID %q did not provide capture group for %q", actual, placeholder.Key)
		}
		result[placeholder.Key] = submatches[groupIndex]
		groupIndex++
	}
	return result, nil
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return map[string]string{}
	}
	output := make(map[string]string, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}
