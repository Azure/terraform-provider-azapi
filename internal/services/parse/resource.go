package parse

import (
	"fmt"
	"log"
	"net/url"
	"strings"

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

func BuildResourceID(name, parentId, resourceType string) (ResourceId, error) {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return ResourceId{}, err
	}
	resourceDef, err := azure.GetResourceDefinition(azureResourceType, apiVersion)
	if err != nil {
		log.Printf("[ERROR] load embedded schema: %+v\n", err)
	}

	if !utils.IsTopLevelResourceType(azureResourceType) {
		parentIdExpectedType := utils.GetParentType(azureResourceType)
		parentIdType := utils.GetResourceType(parentId)
		if !strings.EqualFold(parentIdExpectedType, parentIdType) {
			return ResourceId{}, fmt.Errorf("`parent_id` is invalid, expect id of `%s`", parentIdExpectedType)
		}
		parts := strings.Split(azureResourceType, "/")
		if len(parts) < 2 {
			// impossible to reach here
			return ResourceId{}, fmt.Errorf("`type` and `parent_id` are not matched")
		}
		lastType := parts[len(parts)-1]
		azureResourceId := fmt.Sprintf("%s/%s/%s", parentId, lastType, name)
		return ResourceId{
			AzureResourceId:   azureResourceId,
			ApiVersion:        apiVersion,
			AzureResourceType: azureResourceType,
			Name:              name,
			ParentId:          parentId,
			ResourceDef:       resourceDef,
		}, nil
	}

	scopeTypes := make([]types.ScopeType, 0)
	if resourceDef != nil {
		for _, scope := range resourceDef.ScopeTypes {
			if scope != types.Unknown {
				scopeTypes = append(scopeTypes, scope)
			}
		}
	}
	parentIdScope := utils.GetScopeType(parentId)
	azureResourceId := ""
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
				if parentIdScope == types.ResourceGroup {
					matchedScope = scope
				}
			case types.Unknown:
			}
		}
		if matchedScope == types.Unknown {
			return ResourceId{}, fmt.Errorf("`parent_id` is invalid, expect id of resource whose scope is %v, but got scope %v", scopeTypes, parentIdScope)
		}
	}

	// build azure resource id
	if strings.EqualFold(azureResourceType, "Microsoft.Resources/resourceGroups") {
		azureResourceId = fmt.Sprintf("%s/resourceGroups/%s", parentId, name)
	} else {
		scopeId := parentId
		if parentId == "/" {
			scopeId = ""
		}
		azureResourceId = fmt.Sprintf("%s/providers/%s/%s", scopeId, azureResourceType, name)
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

func NewResourceID(azureResourceId, resourceType string) (ResourceId, error) {
	azureResourceType, _, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return ResourceId{}, err
	}
	if azureResourceType != utils.GetResourceType(azureResourceId) {
		return ResourceId{}, fmt.Errorf("`resource_id` and `type` are not matched")
	}
	name := utils.GetName(azureResourceId)
	parentId := utils.GetParentId(azureResourceId)
	return BuildResourceID(name, parentId, resourceType)
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

// ResourceID parses a Resource ID into an ResourceId struct
func ResourceID(input string) (ResourceId, error) {
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
	id, err := NewResourceID(azureResourceId, fmt.Sprintf("%s@%s", azureResourceType, apiVersion))
	if err != nil {
		return ResourceId{}, err
	}
	return id, nil
}
