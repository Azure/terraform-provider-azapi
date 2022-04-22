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

resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "west europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "mystorageaccount"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azapi_resource" "project" {
  type      = "Microsoft.Migrate/migrateProjects@2020-05-01"
  name      = "myproject"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      publicNetworkAccess     = "Enabled"
      utilityStorageAccountId = azurerm_storage_account.test.id
    }
  })
}

resource "azapi_resource" "solution" {
  type      = "Microsoft.Migrate/migrateProjects/solutions@2018-09-01-preview"
  name      = "mysolution"
  parent_id = azapi_resource.project.id

  body = jsonencode({
    properties = {
      summary = {
        instanceType  = "Servers"
        migratedCount = 0
      }
    }
  })

}
