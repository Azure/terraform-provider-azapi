---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "virtualHubs/routingIntent"
description: |-
  Manages a Virtual Hub Routing Intent.
---

# Microsoft.Network/virtualHubs/routingIntent - Virtual Hub Routing Intent

This article demonstrates how to use `azapi` provider to manage the Virtual Hub Routing Intent resource in Azure.

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

provider "azurerm" {
  features {
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
      addressPrefix        = "10.0.2.0/24"
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

resource "azurerm_firewall" "test" {
  name                = var.resource_name
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  sku_name            = "AZFW_Hub"
  sku_tier            = "Standard"

  virtual_hub {
    virtual_hub_id  = azapi_resource.virtualHub.id
    public_ip_count = 1
  }
}

resource "azapi_resource" "routingIntent" {
  name      = var.resource_name
  type      = "Microsoft.Network/virtualHubs/routingIntent@2022-09-01"
  parent_id = azapi_resource.virtualHub.id

  body = {
    properties = {
      routingPolicies = [
        {
          name = "InternetTraffic"
          destinations = [
            "Internet"
          ]
          nextHop = azurerm_firewall.test.id
        },
        {
          name = "PrivateTrafficPolicy"
          destinations = [
            "PrivateTraffic"
          ]
          nextHop = azurerm_firewall.test.id
        }
      ]
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/virtualHubs/routingIntent@api-version`. The available api-versions for this resource are: [`2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualHubs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/virtualHubs/routingIntent?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualHubs/{resourceName}/routingIntent/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualHubs/{resourceName}/routingIntent/{resourceName}?api-version=2024-07-01
 ```
