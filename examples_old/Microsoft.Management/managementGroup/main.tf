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

// This is an example to demonstrate how to set `parent_id` for tenant scope resource.
// It's recommended to manage management group with `azurerm_management_group`.
resource "azapi_resource" "test" {
  type      = "Microsoft.Management/managementGroups@2021-04-01"
  name      = "myMgmtGroup"
  parent_id = "/"
}
