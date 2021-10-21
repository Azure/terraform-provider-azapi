package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type UserAssignedIdentitiesId struct {
	SubscriptionId           string
	ResourceGroup            string
	UserAssignedIdentityName string
}

func NewUserAssignedIdentitiesID(subscriptionId, resourceGroup, userAssignedIdentityName string) UserAssignedIdentitiesId {
	return UserAssignedIdentitiesId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		UserAssignedIdentityName: userAssignedIdentityName,
	}
}

func (id UserAssignedIdentitiesId) String() string {
	segments := []string{
		fmt.Sprintf("User Assigned Identity Name %q", id.UserAssignedIdentityName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "User Assigned Identities", segmentsStr)
}

func (id UserAssignedIdentitiesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.UserAssignedIdentityName)
}

// UserAssignedIdentitiesID parses a UserAssignedIdentities ID into an UserAssignedIdentitiesId struct
func UserAssignedIdentitiesID(input string) (*UserAssignedIdentitiesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := UserAssignedIdentitiesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.UserAssignedIdentityName, err = id.PopSegment("userAssignedIdentities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
