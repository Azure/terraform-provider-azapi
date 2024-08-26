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
  name                 = "myapp"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  service_plan_id      = azurerm_service_plan.test.id
  storage_account_name = azurerm_storage_account.test.name

  site_config {}
}

resource "azurerm_linux_function_app_slot" "test" {
  name                 = "myslot"
  function_app_id      = azurerm_linux_function_app.test.id
  storage_account_name = azurerm_storage_account.test.name

  site_config {}
}

data "azapi_resource_action" "test" {
  type                   = "Microsoft.Web/sites/slots@2022-03-01"
  resource_id            = azurerm_linux_function_app_slot.test.id
  action                 = "host/default/listkeys"
  response_export_values = ["*"]
}

output "output1" {
  value = data.azapi_resource_action.test.output
}
