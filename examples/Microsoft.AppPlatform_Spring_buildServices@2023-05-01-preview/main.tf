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

data "azapi_resource_id" "buildService" {
  type      = "Microsoft.AppPlatform/Spring/buildServices@2023-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = "default"
}

resource "azapi_resource_action" "buildService" {
  type        = "Microsoft.AppPlatform/Spring/buildServices@2023-05-01-preview"
  resource_id = data.azapi_resource_id.buildService.id
  method      = "PUT"
  body = {
    properties = {
    }
  }
  response_export_values = ["*"]
}

data "azapi_resource" "buildService" {
  type      = "Microsoft.AppPlatform/Spring/buildServices@2023-05-01-preview"
  name      = "default"
  parent_id = azapi_resource.Spring.id

  depends_on = [azapi_resource_action.buildService]
}
