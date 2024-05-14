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

resource "azapi_resource" "iotApp" {
  type      = "Microsoft.IoTCentral/iotApps@2021-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      displayName         = var.resource_name
      publicNetworkAccess = "Enabled"
      subdomain           = "subdomain-2306300333537"
      template            = "iotc-pnp-preview@1.0.0"
    }
    sku = {
      name = "ST1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

