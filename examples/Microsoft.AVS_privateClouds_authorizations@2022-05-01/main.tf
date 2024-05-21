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
  default = "centralus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "privateCloud" {
  type      = "Microsoft.AVS/privateClouds@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      internet = "Disabled"
      managementCluster = {
        clusterSize = 3
      }
      networkBlock = "192.168.48.0/22"
    }
    sku = {
      name = "av36"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorization" {
  type                      = "Microsoft.AVS/privateClouds/authorizations@2022-05-01"
  parent_id                 = azapi_resource.privateCloud.id
  name                      = var.resource_name
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

