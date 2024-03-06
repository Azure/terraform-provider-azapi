package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ActionResource struct{}

func TestAccActionResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_basicWhenDestroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicWhenDestroy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That("azapi_resource_action.test").Key("output").HasValue("{}"),
			),
		},
		{
			Destroy: true,
			Config:  r.basicWhenDestroy(data),
			Check:   resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_registerResourceProvider(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.registerResourceProvider(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_nonstandardLRO(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.nonstandardLRO(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ActionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_resource_action" "list" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = azapi_resource.test.id
  action      = "agentRegistrationInformation/regenerateKey"
  body = jsonencode({
    keyName = "primary"
  })
  depends_on = [
    data.azapi_resource_action.list
  ]
}
`, GenericResource{}.defaultTag(data))
}

func (r ActionResource) basicWhenDestroy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_resource_action" "list" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = azapi_resource.test.id
  when        = "destroy"
  action      = "agentRegistrationInformation/regenerateKey"
  body = jsonencode({
    keyName = "primary"
  })
  depends_on = [
    data.azapi_resource_action.list
  ]
  response_export_values = ["*"]
}
`, GenericResource{}.defaultTag(data))
}

func (r ActionResource) registerResourceProvider() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Resources/providers@2021-04-01"
  resource_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Compute"
  action      = "register"
  method      = "POST"
}
`
}

func (r ActionResource) providerAction() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Cache@2023-04-01"
  resource_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Cache"
  action      = "CheckNameAvailability"
  body = jsonencode({
    type = "Microsoft.Cache/Redis"
    name = "cacheName"
  })
}
`
}

func (r ActionResource) nonstandardLRO(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[2]s"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestsp%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Windows"
  sku_name            = "Y1"
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctestfa%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  service_plan_id            = azurerm_service_plan.test.id

  site_config {}
  ftp_publish_basic_authentication_enabled = false
}

data "azapi_resource_id" "host" {
  type      = "Microsoft.Web/sites/host@2022-03-01"
  parent_id = azurerm_windows_function_app.test.id
  name      = "default"
}

data "azapi_resource_id" "functionKey" {
  type      = "Microsoft.Web/sites/host/functionKeys@2022-03-01"
  parent_id = data.azapi_resource_id.host.id
  name      = "tf_key"
}

resource "azapi_resource_action" "test" {
  type                   = "Microsoft.Web/sites/host/functionKeys@2022-03-01"
  resource_id            = data.azapi_resource_id.functionKey.id
  method                 = "PUT"
  response_export_values = ["*"]
  body = jsonencode({
    properties = {
      name  = "test_key"
      value = "test_value"
    }
  })
}


`, data.LocationPrimary, data.RandomStringOfLength(10))
}
