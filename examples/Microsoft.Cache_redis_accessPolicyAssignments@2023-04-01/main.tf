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
  skip_provider_registration = false
}

data "azurerm_client_config" "test" {
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
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

resource "azapi_resource" "accessPolicyAssignment" {
  type      = "Microsoft.Cache/redis/accessPolicyAssignments@2024-03-01"
  name      = var.resource_name
  parent_id = azapi_resource.redis.id
  body = {
    properties = {
      accessPolicyName = "Data Contributor"
      objectId         = data.azurerm_client_config.test.object_id
      objectIdAlias    = "ServicePrincipal"
    }
  }
}
