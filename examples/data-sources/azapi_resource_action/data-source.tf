terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Basic"
}

data "azapi_resource_action" "example" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azurerm_automation_account.example.id
  action                 = "listKeys"
  response_export_values = ["*"]
}
