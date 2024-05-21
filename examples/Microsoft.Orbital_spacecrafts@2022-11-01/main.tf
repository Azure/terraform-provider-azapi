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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "spacecraft" {
  type      = "Microsoft.Orbital/spacecrafts@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      links = [
        {
          bandwidthMHz       = 100
          centerFrequencyMHz = 101
          direction          = "Uplink"
          name               = "linkname"
          polarization       = "LHCP"
        },
      ]
      noradId   = "12345"
      titleLine = "AQUA"
      tleLine1  = "1 23455U 94089A   97320.90946019  .00000140  00000-0  10191-3 0  2621"
      tleLine2  = "2 23455  99.0090 272.6745 0008546 223.1686 136.8816 14.11711747148495"
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

