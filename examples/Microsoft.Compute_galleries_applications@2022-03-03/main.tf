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

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      description = ""
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "application" {
  type      = "Microsoft.Compute/galleries/applications@2022-03-03"
  parent_id = azapi_resource.gallery.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      supportedOSType = "Linux"
    }

  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

