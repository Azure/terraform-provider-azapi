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

resource "azapi_resource" "routeTable" {
  type      = "Microsoft.Network/routeTables@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableBgpRoutePropagation = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.routes]
  }
}

resource "azapi_resource" "route" {
  type      = "Microsoft.Network/routeTables/routes@2022-09-01"
  parent_id = azapi_resource.routeTable.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.1.0.0/16"
      nextHopType   = "VnetLocal"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

