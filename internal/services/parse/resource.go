package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceId struct {
	AzureResourceId   string
	ApiVersion        string
	AzureResourceType string
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
	azureResourceType := utils.GetResourceType(azureResourceId)
	resourceId := ResourceId{
		AzureResourceId:   azureResourceId,
		AzureResourceType: azureResourceType,
		ApiVersion:        idUrl.Query().Get("api-version"),
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
