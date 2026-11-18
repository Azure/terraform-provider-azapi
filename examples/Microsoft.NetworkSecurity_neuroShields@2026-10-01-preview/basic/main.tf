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
  default = "centraluseuap"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "neuroShield" {
  type                      = "Microsoft.NetworkSecurity/neuroShields@2026-10-01-preview"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false

  body = {
    properties = {
      responseMode      = "Manual"
      severityThreshold = "Medium"
      description       = "Adaptive virtual patching example"
    }
  }
}
