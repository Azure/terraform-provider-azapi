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
      peeringLocation = "Airtel-Chennai2-CLS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "ExpressRoutePort2" {
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
      name   = "Standard_MeteredData"
      tier   = "Standard"
    }
  }
  ignore_casing             = true
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "expressRouteCircuit2" {
  type      = "Microsoft.Network/expressRouteCircuits@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authorizationKey = ""
      bandwidthInGbps  = 5
      expressRoutePort = {
        id = azapi_resource.ExpressRoutePort2.id
      }
    }
    sku = {
      family = "MeteredData"
      name   = "Standard_MeteredData"
      tier   = "Standard"
    }
  }
  ignore_casing             = true
  schema_validation_enabled = false
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
      secondaryPeerAddressPrefix = "192.168.1.0/30"
      sharedKey                  = "ItsASecret"
      state                      = "Enabled"
      vlanId                     = 100
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "peering2" {
  type      = "Microsoft.Network/expressRouteCircuits/peerings@2022-07-01"
  parent_id = azapi_resource.expressRouteCircuit2.id
  name      = "AzurePrivatePeering"
  body = {
    properties = {
      azureASN                   = 12076
      gatewayManagerEtag         = ""
      peerASN                    = 100
      peeringType                = "AzurePrivatePeering"
      primaryPeerAddressPrefix   = "192.168.1.0/30"
      secondaryPeerAddressPrefix = "192.168.1.0/30"
      sharedKey                  = "ItsASecret"
      state                      = "Enabled"
      vlanId                     = 100
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "connection" {
  type      = "Microsoft.Network/expressRouteCircuits/peerings/connections@2022-07-01"
  parent_id = azapi_resource.peering.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "192.169.8.0/29"
      expressRouteCircuitPeering = {
        id = azapi_resource.peering.id
      }
      peerExpressRouteCircuitPeering = {
        id = azapi_resource.peering2.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

