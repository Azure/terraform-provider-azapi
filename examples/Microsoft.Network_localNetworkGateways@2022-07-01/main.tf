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

resource "azapi_resource" "localNetworkGateway" {
  type      = "Microsoft.Network/localNetworkGateways@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      gatewayIpAddress = "168.62.225.23"
      localNetworkAddressSpace = {
        addressPrefixes = [
          "10.1.1.0/24",
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

