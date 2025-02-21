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

resource "azapi_resource" "powerShell72Module" {
  type      = "Microsoft.Automation/automationAccounts/powerShell72Modules@2020-01-13-preview"
  parent_id = azapi_resource.automationAccount.id
  name      = "xActiveDirectory"
  body = {
    properties = {
      contentLink = {
        uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

