terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-rg"
  location = "westus"
  body     = {}
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "exampleaccount"
  location  = "westus"
  body = {
    kind = "StorageV2"
    sku = {
      name = "Premium_LRS"
    }
  }
}

ephemeral "azapi_resource_action" "listKeys" {
  type        = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id = azapi_resource.storageAccount.id
  action      = "listKeys"
  method      = "POST"
  response_export_values = {
    all = "@"
  }
}
