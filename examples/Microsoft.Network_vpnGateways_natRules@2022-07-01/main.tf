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

resource "azapi_resource" "virtualWan" {
  type      = "Microsoft.Network/virtualWans@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allowBranchToBranchTraffic     = true
      disableVpnEncryption           = false
      office365LocalBreakoutCategory = "None"
      type                           = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualHub" {
  type      = "Microsoft.Network/virtualHubs@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressPrefix        = "10.0.0.0/24"
      hubRoutingPreference = "ExpressRoute"
      virtualRouterAutoScaleConfiguration = {
        minCapacity = 2
      }
      virtualWan = {
        id = azapi_resource.virtualWan.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "vpnGateway" {
  type      = "Microsoft.Network/vpnGateways@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enableBgpRouteTranslationForNat = false
      isRoutingPreferenceInternet     = false
      virtualHub = {
        id = azapi_resource.virtualHub.id
      }
      vpnGatewayScaleUnit = 1
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  timeouts {
    create = "180m"
    update = "180m"
    delete = "60m"
  }
}

resource "azapi_resource" "natRule" {
  type      = "Microsoft.Network/vpnGateways/natRules@2022-07-01"
  parent_id = azapi_resource.vpnGateway.id
  name      = var.resource_name
  body = {
    properties = {
      externalMappings = [
        {
          addressSpace = "192.168.21.0/26"
        },
      ]
      internalMappings = [
        {
          addressSpace = "10.4.0.0/26"
        },
      ]
      mode = "EgressSnat"
      type = "Static"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

