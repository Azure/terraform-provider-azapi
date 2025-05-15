---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "virtualNetworks"
description: |-
  Manages a virtual network including any configured subnets.
---

# Microsoft.Network/virtualNetworks - virtual network including any configured subnets

This article demonstrates how to use `azapi` provider to manage the virtual network including any configured subnets resource in Azure.

## Example Usage

### with_ipam_pool

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    time = {
      source = "Hashicorp/time"
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

data "azapi_client_config" "current" {}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description                 = ""
      networkManagerScopeAccesses = []
      networkManagerScopes = {
        managementGroups = []
        subscriptions = [
          "/subscriptions/${data.azapi_client_config.current.subscription_id}"
        ]
      }
    }
  }
  tags = {
    "sampleTag" : "sampleTag"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [azapi_resource.networkManager]

  destroy_duration = "60s"
}

resource "azapi_resource" "ipamPool" {
  type      = "Microsoft.Network/networkManagers/ipamPools@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressPrefixes = [
        "10.0.0.0/24",
      ]
      description    = "Test description."
      parentPoolName = ""
      displayName    = "testDisplayName"
    }
  }
  tags = {
    "myTag" : "testTag"
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false

  depends_on = [time_sleep.wait_60_seconds]
}

resource "azapi_resource" "vnet_withIPAM" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        ipamPoolPrefixAllocations = [
          {
            numberOfIpAddresses = "100",
            pool = {
              id = azapi_resource.ipamPool.id
            }
          }
        ]
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/virtualNetworks@api-version`. The available api-versions for this resource are: [`2015-05-01-preview`, `2015-06-15`, `2016-03-30`, `2016-06-01`, `2016-09-01`, `2016-12-01`, `2017-03-01`, `2017-03-30`, `2017-06-01`, `2017-08-01`, `2017-09-01`, `2017-10-01`, `2017-11-01`, `2018-01-01`, `2018-02-01`, `2018-04-01`, `2018-06-01`, `2018-07-01`, `2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/virtualNetworks?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{resourceName}?api-version=2024-05-01
 ```
