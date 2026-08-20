package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
)

var (
	// ErrInvalidResourceType indicates that a resource type is not in standard AzAPI form.
	ErrInvalidResourceType = errors.New("resource type must be in the form provider/type@api-version")

	// ErrInvalidAPIVersion indicates that a resource type does not include an API version.
	ErrInvalidAPIVersion = errors.New("resource type must include an API version")

	// ErrInvalidProperty indicates that a top-level property name is empty.
	ErrInvalidProperty = errors.New("property name must not be empty")

	// ErrDefinitionUnavailable indicates that no embedded schema definition could be loaded.
	ErrDefinitionUnavailable = errors.New("resource schema definition is unavailable")
)

// PropertyStatus describes whether a top-level property can be written to a resource.
type PropertyStatus string

const (
	// PropertyStatusWritable means the property is present and writable.
	PropertyStatusWritable PropertyStatus = "writable"

	// PropertyStatusReadOnly means the property is present but cannot be written.
	PropertyStatusReadOnly PropertyStatus = "read_only"

	// PropertyStatusUnsupported means the property is absent from the resource definition.
	PropertyStatusUnsupported PropertyStatus = "unsupported"

	// PropertyStatusUnknown means the property could not be inspected. The returned error explains why.
	PropertyStatusUnknown PropertyStatus = "unknown"
)

// TopLevelPropertyStatus returns the capability of property in resourceType.
//
// resourceType must use the standard AzAPI format, such as
// "Microsoft.Resources/resourceGroups@2021-04-01". The result is PropertyStatusUnknown and a
// non-nil error when the input is invalid or the embedded definition is unavailable.
func TopLevelPropertyStatus(resourceType, property string) (PropertyStatus, error) {
	resourceName, apiVersion, err := parseResourceType(resourceType)
	if err != nil {
		return PropertyStatusUnknown, err
	}
	if property == "" || property != strings.TrimSpace(property) {
		return PropertyStatusUnknown, ErrInvalidProperty
	}

	definition, err := azure.GetResourceDefinition(resourceName, apiVersion)
	if err != nil {
		return PropertyStatusUnknown, fmt.Errorf("%w for %q: %v", ErrDefinitionUnavailable, resourceType, err)
	}
	if definition == nil || definition.Body == nil || definition.Body.Type == nil {
		return PropertyStatusUnknown, fmt.Errorf("%w for %q: resource body is unavailable", ErrDefinitionUnavailable, resourceType)
	}

	status, ok := propertyStatus(*definition.Body.Type, property)
	if !ok {
		return PropertyStatusUnknown, fmt.Errorf("%w for %q: resource body is not an object", ErrDefinitionUnavailable, resourceType)
	}
	if definition.IsReadOnly() {
		if status == PropertyStatusUnsupported {
			return status, nil
		}
		return PropertyStatusReadOnly, nil
	}
	return status, nil
}

func parseResourceType(resourceType string) (string, string, error) {
	if resourceType == "" || resourceType != strings.TrimSpace(resourceType) || strings.Count(resourceType, "@") != 1 {
		return "", "", ErrInvalidResourceType
	}

	resourceName, apiVersion, _ := strings.Cut(resourceType, "@")
	if resourceName == "" {
		return "", "", ErrInvalidResourceType
	}
	if apiVersion == "" {
		return "", "", ErrInvalidAPIVersion
	}
	return resourceName, apiVersion, nil
}

func propertyStatus(body types.TypeBase, property string) (PropertyStatus, bool) {
	switch body := body.(type) {
	case *types.ObjectType:
		return objectPropertyStatus(body.Properties, property), true
	case *types.DiscriminatedObjectType:
		return objectPropertyStatus(body.BaseProperties, property), true
	default:
		return PropertyStatusUnknown, false
	}
}

func objectPropertyStatus(properties map[string]types.ObjectProperty, property string) PropertyStatus {
	propertyDefinition, ok := properties[property]
	if !ok {
		return PropertyStatusUnsupported
	}
	if propertyDefinition.IsReadOnly() {
		return PropertyStatusReadOnly
	}
	return PropertyStatusWritable
}
