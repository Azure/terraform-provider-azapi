package parse

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/utils"
)

type ResourceId struct {
	AzureResourceId   string
	ApiVersion        string
	AzureResourceType string
	Name              string
	ParentId          string
	ResourceDef       *types.ResourceType
}

// NewResourceIDWithNestedResourceNames constructs a nested resource ID from the given resource names, parent ID, and resource type.
func NewResourceIDWithNestedResourceNames(resourceNames []string, parentId, resourceType string) (ResourceId, error) {
	if len(resourceNames) == 0 {
		return ResourceId{}, fmt.Errorf("resource names cannot be empty")
	}

	// Append default "@latest" api version if not provided
	resourceType = utils.TryAppendDefaultApiVersion(resourceType)

	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return ResourceId{}, err
	}

	resourceProvider, resourceTypeParts, err := utils.GetAzureResourceTypeParts(azureResourceType)
	if err != nil {
		return ResourceId{}, err
	}

	// Ensure the number of resource names matches the number of resource type parts
	if len(resourceNames) != len(resourceTypeParts) {
		return ResourceId{}, fmt.Errorf("number of resource names does not match the number of resource type parts, expected %d, got %d", len(resourceTypeParts), len(resourceNames))
	}

	currentResourceType := resourceProvider

	// Build resource ID for each nested resource
	for i, resourceTypePart := range resourceTypeParts {
		// Final resource type part
		if i == len(resourceTypeParts)-1 {
			return NewResourceID(resourceNames[i], parentId, resourceType)
		}

		// Intermediate resource type part
		currentResourceType += "/" + resourceTypePart

		parentResourceID, err := NewResourceIDSkipScopeValidation(resourceNames[i], parentId, utils.GetAzureResourceType(currentResourceType, apiVersion))
		if err != nil {
			return ResourceId{}, err
		}

		parentId = parentResourceID.AzureResourceId
	}

	return ResourceId{}, fmt.Errorf("failed to build resource id for nested resources")
}

func NewResourceID(name, parentId, resourceType string) (ResourceId, error) {
	return newResourceID(name, parentId, resourceType, false)
}

func NewResourceIDSkipScopeValidation(name, parentId, resourceType string) (ResourceId, error) {
	return newResourceID(name, parentId, resourceType, true)
}

func newResourceID(name, parentId, resourceType string, skipScopeValidation bool) (ResourceId, error) {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return ResourceId{}, err
	}
	resourceDef, err := azure.GetResourceDefinition(azureResourceType, apiVersion)
	if err != nil {
		log.Printf("[WARN] load embedded schema: %+v\n", err)
	}

	azureResourceId := ""
	switch {
	case !strings.Contains(azureResourceType, "/"):
		// case 0: resource type is a provider type
		// avoid duplicated `/` if parent_id is tenant scope
		scopeId := parentId
		if parentId == "/" {
			scopeId = ""
		}
		azureResourceId = fmt.Sprintf("%s/providers/%s", scopeId, name)
	case utils.IsTopLevelResourceType(azureResourceType):
		// case 1: top level resource, verify parent_id providers correct scope
		if !skipScopeValidation {
			if err = validateParentIdScope(resourceDef, parentId); err != nil {
				return ResourceId{}, fmt.Errorf("`parent_id is invalid`: %+v", err)
			}
		}

		// build azure resource id
		switch azureResourceType {
		case arm.ResourceGroupResourceType.String():
			azureResourceId = fmt.Sprintf("%s/resourceGroups/%s", parentId, name)
		case arm.SubscriptionResourceType.String():
			azureResourceId = fmt.Sprintf("/subscriptions/%s", name)
		case arm.TenantResourceType.String():
			azureResourceId = "/"
		case arm.ProviderResourceType.String():
			// avoid duplicated `/` if parent_id is tenant scope
			scopeId := parentId
			if parentId == "/" {
				scopeId = ""
			}
			azureResourceId = fmt.Sprintf("%s/providers/%s", scopeId, name)
		default:
			// avoid duplicated `/` if parent_id is tenant scope
			scopeId := parentId
			if parentId == "/" {
				scopeId = ""
			}
			azureResourceId = fmt.Sprintf("%s/providers/%s/%s", scopeId, azureResourceType, name)
		}
	default:
		// case 2: child resource, verify parent_id's type matches with resource type's parent type
		if err = validateParentIdType(azureResourceType, parentId); err != nil {
			return ResourceId{}, fmt.Errorf("`parent_id is invalid`: %+v", err)
		}

		// build azure resource id
		lastType := azureResourceType[strings.LastIndex(azureResourceType, "/")+1:]
		azureResourceId = fmt.Sprintf("%s/%s/%s", parentId, lastType, name)
	}

	return ResourceId{
		AzureResourceId:   azureResourceId,
		ApiVersion:        apiVersion,
		AzureResourceType: azureResourceType,
		Name:              name,
		ParentId:          parentId,
		ResourceDef:       resourceDef,
	}, nil
}

// ResourceIDWithResourceType parses a Resource ID and resource type into an ResourceId struct
func ResourceIDWithResourceType(azureResourceId, resourceType string) (ResourceId, error) {
	azureResourceType, _, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return ResourceId{}, err
	}
	resourceTypeFromId := utils.GetResourceType(azureResourceId)
	// if resource type is a provider type, then `type` should be either `Microsoft.Foo` or `Microsoft.Resources/providers`
	if strings.EqualFold(arm.ProviderResourceType.String(), resourceTypeFromId) {
		if strings.Contains(azureResourceType, "/") && !strings.EqualFold(azureResourceType, arm.ProviderResourceType.String()) {
			return ResourceId{}, fmt.Errorf("`resource_id` and `type` are not matched, expect `type` to be a provider type, but got %s", azureResourceType)
		}
	} else {
		if !strings.EqualFold(azureResourceType, resourceTypeFromId) {
			return ResourceId{}, fmt.Errorf("`resource_id` and `type` are not matched, expect `type` to be %s, but got %s", resourceTypeFromId, azureResourceType)
		}
	}

	name := utils.GetName(azureResourceId)
	parentId := utils.GetParentId(azureResourceId)
	id, err := NewResourceID(name, parentId, resourceType)
	if err != nil {
		return ResourceId{}, err
	}
	// The generated resource id is based on the resource type whose case might be different from the input resource id.
	// So we set the generated resource id to the input value.
	id.AzureResourceId = azureResourceId
	return id, nil
}

// ResourceIDWithApiVersion parses a Resource ID which contains api-version into an ResourceId struct
func ResourceIDWithApiVersion(input string) (ResourceId, error) {
	idUrl, err := url.Parse(input)
	if err != nil {
		return ResourceId{}, err
	}

	azureResourceId := idUrl.Path
	apiVersion := idUrl.Query().Get("api-version")

	if azureResourceId == "" {
		return ResourceId{}, fmt.Errorf("ID was missing the 'azure resource id' element")
	}

	if apiVersion == "" {
		return ResourceId{}, fmt.Errorf("ID was missing the 'api-version' element")
	}

	azureResourceType := utils.GetResourceType(azureResourceId)
	id, err := ResourceIDWithResourceType(azureResourceId, fmt.Sprintf("%s@%s", azureResourceType, apiVersion))
	if err != nil {
		return ResourceId{}, err
	}
	return id, nil
}

func (id ResourceId) String() string {
	segments := []string{
		fmt.Sprintf("ResourceId %q", id.AzureResourceId),
		fmt.Sprintf("Api Version %q", id.ApiVersion),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource", segmentsStr)
}

func (id ResourceId) ID() string {
	return id.AzureResourceId
}

func validateParentIdScope(resourceDef *types.ResourceType, parentId string) error {
	if resourceDef != nil {
		scopeTypes := make([]types.ScopeType, 0)
		for _, scope := range resourceDef.ScopeTypes {
			if scope != types.Unknown {
				scopeTypes = append(scopeTypes, scope)
			}
		}

		parentIdScope := utils.GetScopeType(parentId)
		// known scope, use `type` to verify `parent_id`
		if len(scopeTypes) != 0 {
			// check parent_id's scope
			matchedScope := types.Unknown
			for _, scope := range scopeTypes {
				switch scope {
				case types.Tenant, types.ManagementGroup, types.Subscription, types.ResourceGroup:
					if parentIdScope == scope {
						matchedScope = scope
					}
				case types.Extension:
					// skip checking the parent Id's scope because extension resource could be applied to any scope
					matchedScope = scope
				case types.Unknown:
				}
			}
			if matchedScope == types.Unknown {
				return fmt.Errorf("expect ID of resource whose scope is %v, but got scope %v", scopeTypes, parentIdScope)
			}
		}
	}
	return nil
}

func validateParentIdType(azureResourceType string, parentId string) error {
	parentIdExpectedType := utils.GetParentType(azureResourceType)
	parentIdType := utils.GetResourceType(parentId)
	if !strings.EqualFold(parentIdExpectedType, parentIdType) {
		return fmt.Errorf("expect ID of `%s`", parentIdExpectedType)
	}
	return nil
}
