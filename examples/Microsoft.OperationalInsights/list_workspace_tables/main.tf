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

data "azapi_operation" "test" {
  type                   = "Microsoft.OperationalInsights/workspaces@2021-12-01-preview"
  resource_id            = azurerm_log_analytics_workspace.test.id
  operation              = "tables"
  method                 = "GET"
  response_export_values = ["*"]
}

output "first_table_name" {
  value = jsondecode(data.azapi_operation.test.output).value.0.name
}

output "count" {
  value = length(jsondecode(data.azapi_operation.test.output).value)
}
