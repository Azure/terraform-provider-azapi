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

resource "azapi_resource" "expressRouteCircuit" {
  type      = "Microsoft.Network/expressRouteCircuits@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authorizationKey = ""
      serviceProviderProperties = {
        bandwidthInMbps     = 50
        peeringLocation     = "Silicon Valley"
        serviceProviderName = "Equinix"
      }
    }
    sku = {
      family = "MeteredData"
      name   = "Standard_MeteredData"
      tier   = "Standard"
    }
    tags = {
      Environment = "production"
      Purpose     = "AcceptanceTests"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorization" {
  type      = "Microsoft.Network/expressRouteCircuits/authorizations@2022-07-01"
  parent_id = azapi_resource.expressRouteCircuit.id
  name      = var.resource_name
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

