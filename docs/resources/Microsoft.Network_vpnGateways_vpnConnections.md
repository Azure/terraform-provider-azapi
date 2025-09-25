---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "vpnGateways/vpnConnections"
description: |-
  Manages a VPN Gateway Connection.
---

# Microsoft.Network/vpnGateways/vpnConnections - VPN Gateway Connection

This article demonstrates how to use `azapi` provider to manage the VPN Gateway Connection resource in Azure.

## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/vpnGateways/vpnConnections@api-version`. The available api-versions for this resource are: [`2018-04-01`, `2018-06-01`, `2018-07-01`, `2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnGateways/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/vpnGateways/vpnConnections?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnGateways/{resourceName}/vpnConnections/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnGateways/{resourceName}/vpnConnections/{resourceName}?api-version=2024-10-01
 ```
