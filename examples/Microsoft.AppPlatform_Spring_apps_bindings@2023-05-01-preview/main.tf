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
      name = "S0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "app" {
  type      = "Microsoft.AppPlatform/Spring/apps@2023-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      customPersistentDisks = [
      ]
      enableEndToEndTLS = false
      public            = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "redis" {
  type      = "Microsoft.Cache/redis@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        capacity = 2
        family   = "C"
        name     = "Standard"
      }
      enableNonSslPort  = true
      minimumTlsVersion = "1.2"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Cache/redis@2023-04-01"
  resource_id            = azapi_resource.redis.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "binding" {
  type      = "Microsoft.AppPlatform/Spring/apps/bindings@2023-05-01-preview"
  parent_id = azapi_resource.app.id
  name      = var.resource_name
  body = {
    properties = {
      bindingParameters = {
        useSsl = "true"
      }
      key        = data.azapi_resource_action.listKeys.output.primaryKey
      resourceId = azapi_resource.redis.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

