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

data "azurerm_client_config" "current" {
}

// This is an example to demonstrate how to set `parent_id` for subscription scope resource.
// It's recommended to manage resource group with `azurerm_resource_group`.
resource "azapi_resource" "test" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = "myRG"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  location  = "westus"
}
