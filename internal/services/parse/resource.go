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

func NewResourceID(name, parentId, resourceType string) (ResourceId, error) {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return ResourceId{}, err
	}
	resourceDef, err := azure.GetResourceDefinition(azureResourceType, apiVersion)
	if err != nil {
		log.Printf("[WARN] load embedded schema: %+v\n", err)
	}

	azureResourceId := ""
	if utils.IsTopLevelResourceType(azureResourceType) {
		// case 1: top level resource, verify parent_id providers correct scope
		if err = validateParentIdScope(resourceDef, parentId); err != nil {
			return ResourceId{}, fmt.Errorf("`parent_id is invalid`: %+v", err)
		}

		// build azure resource id
		switch azureResourceType {
		case arm.ResourceGroupResourceType.String():
			azureResourceId = fmt.Sprintf("%s/resourceGroups/%s", parentId, name)
		case arm.SubscriptionResourceType.String():
			azureResourceId = fmt.Sprintf("/subscriptions/%s", name)
		default:
			// avoid duplicated `/` if parent_id is tenant scope
			scopeId := parentId
			if parentId == "/" {
				scopeId = ""
			}
			azureResourceId = fmt.Sprintf("%s/providers/%s/%s", scopeId, azureResourceType, name)
		}
	} else {
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
	if resourceTypeFromId := utils.GetResourceType(azureResourceId); azureResourceType != resourceTypeFromId {
		return ResourceId{}, fmt.Errorf("`resource_id` and `type` are not matched, expect `type` to be %s, but got %s", resourceTypeFromId, azureResourceType)
	}
	name := utils.GetName(azureResourceId)
	parentId := utils.GetParentId(azureResourceId)
	return NewResourceID(name, parentId, resourceType)
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
					// only supports extension on a resource group scope resource
					if parentIdScope == types.ResourceGroup {
						matchedScope = scope
					}
				case types.Unknown:
				}
			}
			if matchedScope == types.Unknown {
				return fmt.Errorf("expect id of resource whose scope is %v, but got scope %v", scopeTypes, parentIdScope)
			}
		}
	}
	return nil
}

func validateParentIdType(azureResourceType string, parentId string) error {
	parentIdExpectedType := utils.GetParentType(azureResourceType)
	parentIdType := utils.GetResourceType(parentId)
	if parentIdExpectedType != parentIdType {
		return fmt.Errorf("expect id of `%s`", parentIdExpectedType)
	}
	return nil
}
