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
  body = {
    properties = {
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2022-09-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

data "azapi_resource" "fileService" {
  type      = "Microsoft.Storage/storageAccounts/fileServices@2022-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
}

resource "azapi_resource" "share" {
  type      = "Microsoft.Storage/storageAccounts/fileServices/shares@2022-09-01"
  parent_id = data.azapi_resource.fileService.id
  name      = var.resource_name
  body = {
    properties = {
      accessTier = "Cool"
    }
  }
  response_export_values = ["*"]
}

