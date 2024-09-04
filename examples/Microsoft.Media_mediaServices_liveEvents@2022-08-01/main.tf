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
    kind = "StorageV2"
    sku = {
      name = "Standard_GRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "mediaService" {
  type      = "Microsoft.Media/mediaServices@2021-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      storageAccounts = [
        {
          id   = azapi_resource.storageAccount.id
          type = "Primary"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "liveEvent" {
  type      = "Microsoft.Media/mediaServices/liveEvents@2022-08-01"
  parent_id = azapi_resource.mediaService.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      input = {
        accessControl = {
          ip = {
            allow = [
              {
                address            = "0.0.0.0"
                name               = "AllowAll"
                subnetPrefixLength = 0
              },
            ]
          }
        }
        keyFrameIntervalDuration = "PT6S"
        streamingProtocol        = "RTMP"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

