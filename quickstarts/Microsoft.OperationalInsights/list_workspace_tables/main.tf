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

resource "azurerm_log_analytics_workspace" "test" {
  name                = "myworkspace"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

data "azapi_resource_action" "test" {
  type                   = "Microsoft.OperationalInsights/workspaces@2021-12-01-preview"
  resource_id            = azurerm_log_analytics_workspace.test.id
  action                 = "tables"
  method                 = "GET"
  response_export_values = ["*"]
}

output "first_table_name" {
  value = data.azapi_resource_action.test.output.value.0.name
}

output "count" {
  value = length(data.azapi_resource_action.test.output.value)
}
