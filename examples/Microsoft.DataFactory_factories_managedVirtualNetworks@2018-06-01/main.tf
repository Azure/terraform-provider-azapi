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
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      globalParameters = {
      }
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "managedVirtualNetwork" {
  type      = "Microsoft.DataFactory/factories/managedVirtualNetworks@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = "default"
  body = jsonencode({
    properties = {
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

