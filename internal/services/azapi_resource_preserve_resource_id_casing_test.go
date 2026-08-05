package services_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// serverFarmsCasing matches the casing azurerm writes into state for an App Service Plan.
var serverFarmsCasing = regexp.MustCompile(`/Microsoft\.Web/serverFarms/`)

// TestAccGenericResource_preserveResourceIDCasingAzurermMigration covers the migration
// reported in GH-1120: azurerm_service_plan stores its id as `.../Microsoft.Web/serverFarms/...`,
// the resource is then adopted by azapi_resource whose type is `Microsoft.Web/serverfarms`,
// and a later update rebuilds the id from the type, flipping the casing in state and identity.
func TestAccGenericResource_preserveResourceIDCasingAzurermMigration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	const servicePlan = "azapi_resource.servicePlan"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.azurermServicePlan(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestMatchResourceAttr("azurerm_service_plan.test", "id", serverFarmsCasing),
			),
		},
		{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.azurermServicePlanAdoptedByAzapi(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(servicePlan).ExistsInAzure(r),
				resource.TestMatchResourceAttr(servicePlan, "id", serverFarmsCasing),
			),
		},
		{
			// Forces an update so CreateUpdate recomputes the id from the type.
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.azurermServicePlanAdoptedByAzapi(data, "acctest"),
			Check: resource.ComposeTestCheckFunc(
				check.That(servicePlan).ExistsInAzure(r),
				resource.TestMatchResourceAttr(servicePlan, "id", serverFarmsCasing),
			),
		},
	})
}

func (r GenericResource) azurermServicePlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg%[1]s"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestsp%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Windows"
  sku_name            = "P1v3"
}
`, data.RandomString, data.LocationPrimary)
}

func (r GenericResource) azurermServicePlanAdoptedByAzapi(data acceptance.TestData, tag string) string {
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
  from = azurerm_service_plan.test
  lifecycle {
    destroy = false
  }
}

locals {
  // Matches the id azurerm wrote to state for azurerm_service_plan.test.
  service_plan_id = "${azurerm_resource_group.test.id}/providers/Microsoft.Web/serverFarms/acctestsp%[1]s"
}

import {
  to = azapi_resource.servicePlan
  identity = {
    id = local.service_plan_id
  }
}

resource "azapi_resource" "servicePlan" {
  type      = "Microsoft.Web/serverfarms@2023-12-01"
  name      = "acctestsp%[1]s"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = {
    properties = {
      hyperV         = false
      perSiteScaling = false
      reserved       = false
      zoneRedundant  = false
    }
    sku = {
      name = "P1v3"
    }
  }
%[3]s
}
`, data.RandomString, data.LocationPrimary, tags)
}
