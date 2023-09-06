terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    kind = "StorageV2"
    properties = {
    }
    sku = {
      name = "Standard_LRS"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "defenderForStorageSetting" {
  type      = "Microsoft.Security/defenderForStorageSettings@2022-12-01-preview"
  parent_id = azapi_resource.storageAccount.id
  name      = "current"
  body = jsonencode({
    properties = {
      isEnabled = true
      malwareScanning = {
        onUpload = {
          capGBPerMonth = 5000
          isEnabled     = true
        }
      }
      sensitiveDataDiscovery = {
        isEnabled = true
      }
      overrideSubscriptionLevelSettings = true
    }
  })
}
