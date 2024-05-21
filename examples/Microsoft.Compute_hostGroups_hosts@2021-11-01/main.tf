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

resource "azapi_resource" "hostGroup" {
  type      = "Microsoft.Compute/hostGroups@2021-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      platformFaultDomainCount = 2
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "host" {
  type      = "Microsoft.Compute/hostGroups/hosts@2021-11-01"
  parent_id = azapi_resource.hostGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoReplaceOnFailure = true
      licenseType          = "None"
      platformFaultDomain  = 1
    }
    sku = {
      name = "DSv3-Type1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
