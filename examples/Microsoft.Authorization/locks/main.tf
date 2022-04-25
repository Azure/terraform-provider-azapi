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

// This is an example to demonstrate how to set `parent_id` for extension scope resource.
resource "azapi_resource" "lock" {
  type      = "Microsoft.Authorization/locks@2015-01-01"
  name      = "myLock"
  parent_id = azurerm_resource_group.test.id

  body = jsonencode({
    properties = {
      level = "CanNotDelete"
    }
  })
}
