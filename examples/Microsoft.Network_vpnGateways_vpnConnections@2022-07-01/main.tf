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

resource "azapi_resource" "vpnSite" {
  type      = "Microsoft.Network/vpnSites@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.1.0/24",
        ]
      }
      virtualWan = {
        id = azapi_resource.virtualWan.id
      }
      vpnSiteLinks = [
        {
          name = "link1"
          properties = {
            fqdn      = ""
            ipAddress = "10.0.1.1"
            linkProperties = {
              linkProviderName = ""
              linkSpeedInMbps  = 0
            }
          }
        },
        {
          name = "link2"
          properties = {
            fqdn      = ""
            ipAddress = "10.0.1.2"
            linkProperties = {
              linkProviderName = ""
              linkSpeedInMbps  = 0
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "link1" {
  type      = "Microsoft.Network/vpnSites/vpnSiteLinks@2022-07-01"
  parent_id = azapi_resource.vpnSite.id
  name      = "link1"
}

data "azapi_resource_id" "link2" {
  type      = "Microsoft.Network/vpnSites/vpnSiteLinks@2022-07-01"
  parent_id = azapi_resource.vpnSite.id
  name      = "link2"
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

resource "azapi_resource" "vpnConnection" {
  type      = "Microsoft.Network/vpnGateways/vpnConnections@2022-07-01"
  parent_id = azapi_resource.vpnGateway.id
  name      = var.resource_name
  body = {
    properties = {
      enableInternetSecurity = false
      remoteVpnSite = {
        id = azapi_resource.vpnSite.id
      }
      vpnLinkConnections = [
        {
          name = "link1"
          properties = {
            connectionBandwidth            = 10
            enableBgp                      = false
            enableRateLimiting             = false
            routingWeight                  = 0
            useLocalAzureIpAddress         = false
            usePolicyBasedTrafficSelectors = false
            vpnConnectionProtocolType      = "IKEv2"
            vpnGatewayCustomBgpAddresses = [
            ]
            vpnLinkConnectionMode = "Default"
            vpnSiteLink = {
              id = data.azapi_resource_id.link1.id
            }
          }
        },
        {
          name = "link2"
          properties = {
            connectionBandwidth            = 10
            enableBgp                      = false
            enableRateLimiting             = false
            routingWeight                  = 0
            useLocalAzureIpAddress         = false
            usePolicyBasedTrafficSelectors = false
            vpnConnectionProtocolType      = "IKEv2"
            vpnGatewayCustomBgpAddresses = [
            ]
            vpnLinkConnectionMode = "Default"
            vpnSiteLink = {
              id = data.azapi_resource_id.link2.id
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

