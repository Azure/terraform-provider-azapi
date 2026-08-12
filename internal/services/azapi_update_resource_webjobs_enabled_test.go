package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericUpdateResource_webJobsEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}
	config := r.webJobsEnabled(data)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            config,
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That("data.azapi_resource.web_app").Key("output.properties.siteConfig.webJobsEnabled").HasValue("true"),
			),
		},
		{
			Config:            config,
			ExternalProviders: externalProvidersAzurerm(),
			PlanOnly:          true,
		},
	})
}

func (GenericUpdateResource) webJobsEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]s"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestsp-%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctest-webjobs-%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Web/sites@2025-03-01"
  resource_id = azurerm_linux_web_app.test.id
  body = {
    properties = {
      siteConfig = {
        webJobsEnabled = true
      }
    }
  }
}

data "azapi_resource" "web_app" {
  type                   = "Microsoft.Web/sites@2025-03-01"
  resource_id            = azurerm_linux_web_app.test.id
  response_export_values = ["properties.siteConfig.webJobsEnabled"]

  depends_on = [azapi_update_resource.test]
}
`, data.RandomString, data.LocationPrimary)
}
