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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "appliance" {
  type      = "Microsoft.ResourceConnector/appliances@2022-10-27"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-appliance"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      distro = "AKSEdge"
      infrastructureConfig = {
        provider = "VMWare"
      }
    }
  }
}

