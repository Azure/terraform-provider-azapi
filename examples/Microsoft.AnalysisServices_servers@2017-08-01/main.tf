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

resource "azapi_resource" "server" {
  type      = "Microsoft.AnalysisServices/servers@2017-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      asAdministrators = {
        members = [
        ]
      }
      ipV4FirewallSettings = {
        enablePowerBIService = false
        firewallRules = [
        ]
      }
    }
    sku = {
      name = "B1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

