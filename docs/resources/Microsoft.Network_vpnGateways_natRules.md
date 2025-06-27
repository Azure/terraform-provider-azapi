---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "vpnGateways/natRules"
description: |-
  Manages a VPN Gateway NAT Rule.
---

# Microsoft.Network/vpnGateways/natRules - VPN Gateway NAT Rule

This article demonstrates how to use `azapi` provider to manage the VPN Gateway NAT Rule resource in Azure.

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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/vpnGateways/natRules@api-version`. The available api-versions for this resource are: [`2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnGateways/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/vpnGateways/natRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnGateways/{resourceName}/natRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnGateways/{resourceName}/natRules/{resourceName}?api-version=2024-07-01
 ```
