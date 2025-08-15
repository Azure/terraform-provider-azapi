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

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devCenters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned"
  }
  body = {
    properties = {}
  }
}

resource "azapi_resource" "project" {
  type      = "Microsoft.DevCenter/projects@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-proj"
  location  = var.location
  body = {
    properties = {
      description        = ""
      devCenterId        = azapi_resource.devCenter.id
      maxDevBoxesPerUser = 0
    }
  }
}

