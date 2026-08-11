package services_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// userAssignedIdentitiesCasing matches the casing azurerm writes into state for a user-assigned identity.
var userAssignedIdentitiesCasing = regexp.MustCompile(`/Microsoft\.ManagedIdentity/userAssignedIdentities/`)

// TestAccGenericResource_preserveResourceIDCasingAzurermMigration covers the migration
// reported in GH-1120: an azurerm resource stores its id using the API's canonical casing,
// the resource is then adopted by azapi_resource whose type uses different casing,
// and a later update rebuilds the id from the type, flipping the casing in state and identity.
func TestAccGenericResource_preserveResourceIDCasingAzurermMigration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	const userAssignedIdentity = "azapi_resource.userAssignedIdentity"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.azurermUserAssignedIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestMatchResourceAttr("azurerm_user_assigned_identity.test", "id", userAssignedIdentitiesCasing),
			),
		},
		{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.azurermUserAssignedIdentityAdoptedByAzapi(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(userAssignedIdentity).ExistsInAzure(r),
				resource.TestMatchResourceAttr(userAssignedIdentity, "id", userAssignedIdentitiesCasing),
			),
		},
		{
			// Forces an update so CreateUpdate recomputes the id from the type.
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.azurermUserAssignedIdentityAdoptedByAzapi(data, "acctest"),
			Check: resource.ComposeTestCheckFunc(
				check.That(userAssignedIdentity).ExistsInAzure(r),
				resource.TestMatchResourceAttr(userAssignedIdentity, "id", userAssignedIdentitiesCasing),
			),
		},
	})
}

func (r GenericResource) azurermUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg%[1]s"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestid%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomString, data.LocationPrimary)
}

func (r GenericResource) azurermUserAssignedIdentityAdoptedByAzapi(data acceptance.TestData, tag string) string {
	tags := ""
	if tag != "" {
		tags = fmt.Sprintf(`
  tags = {
    environment = %[1]q
  }
`, tag)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azapi" {
  preserve_resource_id_casing = true
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg%[1]s"
  location = "%[2]s"
}

removed {
  from = azurerm_user_assigned_identity.test
  lifecycle {
    destroy = false
  }
}

locals {
  // Matches the id azurerm wrote to state for azurerm_user_assigned_identity.test.
  user_assigned_identity_id = "${azurerm_resource_group.test.id}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctestid%[1]s"
}

import {
  to = azapi_resource.userAssignedIdentity
  identity = {
    id = local.user_assigned_identity_id
  }
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userassignedidentities@2023-01-31"
  name      = "acctestid%[1]s"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
%[3]s
}
`, data.RandomString, data.LocationPrimary, tags)
}
