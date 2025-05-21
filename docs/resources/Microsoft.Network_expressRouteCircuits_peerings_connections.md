---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "expressRouteCircuits/peerings/connections"
description: |-
  Manages a Express Route Circuit Connection.
---

# Microsoft.Network/expressRouteCircuits/peerings/connections - Express Route Circuit Connection

This article demonstrates how to use `azapi` provider to manage the Express Route Circuit Connection resource in Azure.

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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/expressRouteCircuits/peerings/connections@api-version`. The available api-versions for this resource are: [`2018-02-01`, `2018-04-01`, `2018-06-01`, `2018-07-01`, `2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{resourceName}/peerings/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/expressRouteCircuits/peerings/connections?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{resourceName}/peerings/{resourceName}/connections/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteCircuits/{resourceName}/peerings/{resourceName}/connections/{resourceName}?api-version=2024-05-01
 ```
