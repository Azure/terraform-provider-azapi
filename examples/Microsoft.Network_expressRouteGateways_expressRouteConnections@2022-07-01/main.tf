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

resource "azapi_resource" "ExpressRoutePort" {
  type      = "Microsoft.Network/ExpressRoutePorts@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      bandwidthInGbps = 10
      encapsulation   = "Dot1Q"
      peeringLocation = "CDC-Canberra"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
      addressPrefix        = "10.0.1.0/24"
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

resource "azapi_resource" "expressRouteCircuit" {
  type      = "Microsoft.Network/expressRouteCircuits@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authorizationKey = ""
      bandwidthInGbps  = 5
      expressRoutePort = {
        id = azapi_resource.ExpressRoutePort.id
      }
    }
    sku = {
      family = "MeteredData"
      name   = "Premium_MeteredData"
      tier   = "Premium"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}

resource "azapi_resource" "peering" {
  type      = "Microsoft.Network/expressRouteCircuits/peerings@2022-07-01"
  parent_id = azapi_resource.expressRouteCircuit.id
  name      = "AzurePrivatePeering"
  body = {
    properties = {
      azureASN                   = 12076
      gatewayManagerEtag         = ""
      peerASN                    = 100
      peeringType                = "AzurePrivatePeering"
      primaryPeerAddressPrefix   = "192.168.1.0/30"
      secondaryPeerAddressPrefix = "192.168.2.0/30"
      sharedKey                  = "ItsASecret"
      state                      = "Enabled"
      vlanId                     = 100
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "expressRouteGateway" {
  type      = "Microsoft.Network/expressRouteGateways@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allowNonVirtualWanTraffic = false
      autoScaleConfiguration = {
        bounds = {
          min = 1
        }
      }
      virtualHub = {
        id = azapi_resource.virtualHub.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "expressRouteConnection" {
  type      = "Microsoft.Network/expressRouteGateways/expressRouteConnections@2022-07-01"
  parent_id = azapi_resource.expressRouteGateway.id
  name      = var.resource_name
  body = {
    properties = {
      enableInternetSecurity = false
      expressRouteCircuitPeering = {
        id = azapi_resource.peering.id
      }
      expressRouteGatewayBypass = false
      routingConfiguration = {
      }
      routingWeight = 0
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

