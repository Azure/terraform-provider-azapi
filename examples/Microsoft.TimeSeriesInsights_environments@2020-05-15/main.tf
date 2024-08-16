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

resource "azapi_resource" "environment" {
  type      = "Microsoft.TimeSeriesInsights/environments@2020-05-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "Gen1"
    properties = {
      dataRetentionTime            = "P30D"
      storageLimitExceededBehavior = "PurgeOldData"
    }
    sku = {
      capacity = 1
      name     = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

