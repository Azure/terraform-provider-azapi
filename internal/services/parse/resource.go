package parse

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
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
	parts := strings.Split(resourceType, "@")
	apiVersion := ""
	azureResourceType := ""
	if len(parts) == 2 {
		apiVersion = parts[1]
		azureResourceType = parts[0]
	}

	azureResourceId := ""
	parts = strings.Split(azureResourceType, "/")
	switch {
	case len(parts) <= 1:
		break
	case len(parts) == 2:
		azureResourceId = fmt.Sprintf("%s/providers/%s/%s", parentId, azureResourceType, name)
	default:
		lastType := parts[len(parts)-1]
		azureResourceId = fmt.Sprintf("%s/%s/%s", parentId, lastType, name)
	}

	if len(parentId) != 0 {
		parentType := utils.GetParentType(azureResourceType)
		if len(parentType) != 0 && !strings.EqualFold(parentType, utils.GetResourceType(parentId)) {
			return ResourceId{}, fmt.Errorf("`parent_id` is invalid, expect id of `%s`", parentType)
		}
	}

	resourceDef, err := azure.GetResourceDefinition(azureResourceType, apiVersion)
	if err != nil {
		log.Printf("[ERROR] load embedded schema: %+v\n", err)
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

func NewResourceID(azureResourceId, resourceType string) ResourceId {
	parts := strings.Split(resourceType, "@")
	apiVersion := ""
	azureResourceType := ""
	if len(parts) == 2 {
		apiVersion = parts[1]
		azureResourceType = parts[0]
	}
	return ResourceId{
		AzureResourceId:   azureResourceId,
		ApiVersion:        apiVersion,
		AzureResourceType: azureResourceType,
		Name:              utils.GetName(azureResourceId),
		ParentId:          utils.GetParentId(azureResourceId),
	}
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
	fmtString := "%s?api-version=%s"
	return fmt.Sprintf(fmtString, id.AzureResourceId, id.ApiVersion)
}

// ResourceID parses a Resource ID into an ResourceId struct
func ResourceID(input string) (*ResourceId, error) {
	idUrl, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	azureResourceId := idUrl.Path
	apiVersion := idUrl.Query().Get("api-version")
	azureResourceType := utils.GetResourceType(azureResourceId)
	resourceDef, err := azure.GetResourceDefinition(azureResourceType, apiVersion)
	if err != nil {
		log.Printf("[ERROR] load embedded schema: %+v\n", err)
	}

	resourceId := ResourceId{
		AzureResourceId:   azureResourceId,
		AzureResourceType: azureResourceType,
		ApiVersion:        apiVersion,
		Name:              utils.GetName(azureResourceId),
		ParentId:          utils.GetParentId(azureResourceId),
		ResourceDef:       resourceDef,
	}

	if resourceId.AzureResourceId == "" {
		return nil, fmt.Errorf("ID was missing the 'azure resource id' element")
	}

	if resourceId.ApiVersion == "" {
		return nil, fmt.Errorf("ID was missing the 'api-version' element")
	}

	id, err := resourceids.ParseAzureResourceID(resourceId.AzureResourceId)
	if err != nil {
		return nil, err
	}

	if id.SubscriptionID == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if id.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	return &resourceId, nil
}
