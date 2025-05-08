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

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2021-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "mradministrator"
      administratorLoginPassword    = "thisIsDog11"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource_action" "automaticTuning" {
  resource_id = "${azapi_resource.server.id}/automaticTuning/current"
  type        = "Microsoft.Sql/servers/automaticTuning@2021-11-01"
  method      = "PATCH"
  body = {
    properties = {
      desiredState = "Auto"
      options = {
        forceLastGoodPlan = { desiredState = "Default" }
        createIndex       = { desiredState = "On" }
        dropIndex         = { desiredState = "Off" }
      }
    }
  }
}
