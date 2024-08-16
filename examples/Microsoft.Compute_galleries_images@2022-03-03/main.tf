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

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "image" {
  type      = "Microsoft.Compute/galleries/images@2022-03-03"
  parent_id = azapi_resource.gallery.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      architecture = "x64"
      description  = ""
      disallowed = {
        diskTypes = [
        ]
      }
      features         = null
      hyperVGeneration = "V1"
      identifier = {
        offer     = "AccTesOffer230630032848825313"
        publisher = "AccTesPublisher230630032848825313"
        sku       = "AccTesSku230630032848825313"
      }
      osState             = "Generalized"
      osType              = "Linux"
      privacyStatementUri = ""
      recommended = {
        memory = {
        }
        vCPUs = {
        }
      }
      releaseNoteUri = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

