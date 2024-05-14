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

resource "azapi_resource" "dnsZone" {
  type                      = "Microsoft.Network/dnsZones@2018-05-01"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}.com"
  location                  = "global"
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "AAAA" {
  type      = "Microsoft.Network/dnsZones/AAAA@2018-05-01"
  parent_id = azapi_resource.dnsZone.id
  name      = var.resource_name
  body = {
    properties = {
      AAAARecords = [
        {
          ipv6Address = "2607:f8b0:4009:1803::1005"
        },
        {
          ipv6Address = "2607:f8b0:4009:1803::1006"
        },
      ]
      TTL = 300
      metadata = {
      }
      targetResource = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

