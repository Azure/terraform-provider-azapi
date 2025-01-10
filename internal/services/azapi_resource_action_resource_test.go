package services_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ActionResource struct{}

func (r ActionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	out := false
	return &out, nil
}

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
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicWhenDestroy(data),
			Check:  resource.ComposeTestCheckFunc(),
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
			Config: r.registerResourceProvider(os.Getenv("ARM_SUBSCRIPTION_ID")),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_nonstandardLRO(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config:            r.nonstandardLRO(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check:             resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.timeouts(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.headers(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_queryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.queryParameters(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionResource_sensitiveOutput(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.sensitiveOutput(data),
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
  body = {
    keyName = "primary"
  }
  depends_on = [
    data.azapi_resource_action.list
  ]
}
`, GenericResource{}.identityNone(data))
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
  body = {
    keyName = "primary"
  }
  depends_on = [
    data.azapi_resource_action.list
  ]
  response_export_values = ["*"]
}
`, GenericResource{}.identityNone(data))
}

func (r ActionResource) registerResourceProvider(subscriptionId string) string {
	return fmt.Sprintf(`
resource "azapi_resource_action" "test" {
  type                   = "Microsoft.Resources/providers@2021-04-01"
  resource_id            = "/subscriptions/%s/providers/Microsoft.Compute"
  action                 = "register"
  method                 = "POST"
  response_export_values = ["*"]
}
`, subscriptionId)
}

func (r ActionResource) providerAction(data acceptance.TestData) string {
	return fmt.Sprintf(`

data "azapi_client_config" "current" {}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Cache@2023-04-01"
  resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Cache"
  action      = "CheckNameAvailability"
  body = {
    type = "Microsoft.Cache/Redis"
    name = "%s"
  }
}
`, data.RandomString)
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
  name                            = "acctestsa%[2]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
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

  ftp_publish_basic_authentication_enabled       = false
  webdeploy_publish_basic_authentication_enabled = false
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
  body = {
    properties = {
      name  = "test_key"
      value = "test_value"
    }
  }
}
`, data.LocationPrimary, data.RandomString)
}

func (r ActionResource) timeouts(data acceptance.TestData) string {
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
  body = {
    keyName = "primary"
  }
  depends_on = [
    data.azapi_resource_action.list
  ]
  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
    read   = "10m"
  }
}
`, GenericResource{}.identityNone(data))
}

func (r ActionResource) oldConfig(data acceptance.TestData, subscriptionId string) string {
	return fmt.Sprintf(`

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Cache@2023-04-01"
  resource_id = "/subscriptions/%s/providers/Microsoft.Cache"
  action      = "CheckNameAvailability"
  body = jsonencode({
    type = "Microsoft.Cache/Redis"
    name = "%s"
  })
  response_export_values = ["*"]
}
`, subscriptionId, data.RandomString)
}

func (r ActionResource) headers(data acceptance.TestData) string {
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
  body = {
    keyName = "primary"
  }
  headers = {
    "header1" = "value1"
  }
  depends_on = [
    data.azapi_resource_action.list
  ]
}
`, GenericResource{}.identityNone(data))
}

func (r ActionResource) queryParameters() string {
	return `
data "azapi_client_config" "current" {}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Authorization@2021-06-01"
  resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization"
  action      = "policyDefinitions"
  method      = "GET"
  query_parameters = {
    "$filter" = ["policyType eq 'BuiltIn'"]
  }
  response_export_values = ["*"]
}`
}

func (r ActionResource) sensitiveOutput(data acceptance.TestData) string {
	return fmt.Sprintf(`

data "azapi_client_config" "current" {}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Cache@2023-04-01"
  resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Cache"
  action      = "CheckNameAvailability"
  body = {
    type = "Microsoft.Cache/Redis"
    name = "%s"
  }
  sensitive_response_export_values = ["*"]
}
`, data.RandomString)
}
