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
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "integrationRuntime" {
  type      = "Microsoft.DataFactory/factories/integrationRuntimes@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      description = ""
      type        = "SelfHosted"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

