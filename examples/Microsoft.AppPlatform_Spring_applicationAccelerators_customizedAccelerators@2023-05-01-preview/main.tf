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

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2023-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "E0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "applicationAccelerator" {
  type                      = "Microsoft.AppPlatform/Spring/applicationAccelerators@2023-05-01-preview"
  parent_id                 = azapi_resource.Spring.id
  name                      = "default"
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "customizedAccelerator" {
  type      = "Microsoft.AppPlatform/Spring/applicationAccelerators/customizedAccelerators@2023-05-01-preview"
  parent_id = azapi_resource.applicationAccelerator.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      displayName = ""
      gitRepository = {
        authSetting = {
          authType = "Public"
        }
        branch = "master"
        commit = ""
        gitTag = ""
        url    = "https://github.com/Azure-Samples/piggymetrics"
      }
      iconUrl = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

