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

data "azapi_client_config" "current" {}

variable "resource_name" {
  type = string
}

variable "location" {
  type = string
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

resource "azapi_resource" "environmentType" {
  type      = "Microsoft.DevCenter/devCenters/environmentTypes@2025-02-01"
  parent_id = azapi_resource.devCenter.id
  name      = "${var.resource_name}-envtype"
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

resource "azapi_resource" "environmentType_1" {
  type      = "Microsoft.DevCenter/projects/environmentTypes@2025-02-01"
  parent_id = azapi_resource.project.id
  name      = azapi_resource.environmentType.name
  body = {
    properties = {
      deploymentTargetId = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
      status             = "Enabled"
    }
  }
}

