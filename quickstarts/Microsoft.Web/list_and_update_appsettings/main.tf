terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "myaccount"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "myplan"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

resource "azurerm_linux_function_app" "test" {
  name                = "myapp"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage_account_name = azurerm_storage_account.test.name
  service_plan_id      = azurerm_service_plan.test.id

  site_config {}

  lifecycle {
    // appsettings is already supported in azurerm, this example demostrates how to use azapi_resource_action to update the settings
    ignore_changes = [app_settings]
  }
}

data "azapi_resource" "appsettings" {
  type                   = "Microsoft.Web/sites/config@2022-03-01"
  parent_id              = azurerm_linux_function_app.test.id
  name                   = "appsettings"
  response_export_values = ["*"]
}

output "o1" {
  // appsettings can't be fetched with azapi_resource data source directly
  value = data.azapi_resource.appsettings.output
}

data "azapi_resource_action" "list" {
  type                   = "Microsoft.Web/sites/config@2022-03-01"
  resource_id            = data.azapi_resource.appsettings.id
  action                 = "list"
  method                 = "POST"
  response_export_values = ["*"]
}

output "o2" {
  // appsettings can only be fetched with list action
  value = data.azapi_resource_action.list.output
}

resource "azapi_resource_action" "update" {
  type        = "Microsoft.Web/sites/config@2022-03-01"
  resource_id = data.azapi_resource.appsettings.id
  method      = "PUT"
  body = {
    name = "appsettings"
    // use merge function to combine new settings with existing ones
    properties = merge(
      data.azapi_resource_action.list.output.properties,
      {
        WEBSITES_ENABLE_APP_SERVICE_STORAGE = "false"
      }
    )
  }
  response_export_values = ["*"]
}

output "o3" {
  value = azapi_resource_action.update.output
}
