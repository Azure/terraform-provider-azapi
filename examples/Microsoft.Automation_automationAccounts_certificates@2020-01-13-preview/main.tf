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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "certificate" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2020-01-13-preview"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      base64Value  = "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"
      description  = ""
      isExportable = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

