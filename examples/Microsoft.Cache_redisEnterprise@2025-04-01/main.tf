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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "redisEnterprise" {
  type      = "Microsoft.Cache/redisEnterprise@2025-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      minimumTlsVersion = "1.2"
      encryption        = {}
      highAvailability  = "Enabled"
    }
    sku = {
      name = "Balanced_B0"
    }
  }

  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

