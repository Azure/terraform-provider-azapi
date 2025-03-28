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

provider "azurerm" {
  features {
  }
}

data "azurerm_client_config" "current" {
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

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devcenters@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "SystemAssigned"
      userAssignedIdentities = null
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "devBoxDefinition" {
  type      = "Microsoft.DevCenter/devcenters/devboxdefinitions@2024-10-01-preview"
  parent_id = azapi_resource.devCenter.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hibernateSupport = "Enabled"
      imageReference = {
        id = "${azapi_resource.devCenter.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
      },
      sku = {
        name = "general_i_8c32gb256ssd_v2"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
