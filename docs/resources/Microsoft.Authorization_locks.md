---
subcategory: "Microsoft.Authorization - Azure Resource Manager"
page_title: "locks"
description: |-
  Manages a Management Lock which is scoped to a Subscription, Resource Group or Resource.
---

# Microsoft.Authorization/locks - Management Lock which is scoped to a Subscription, Resource Group or Resource

This article demonstrates how to use `azapi` provider to manage the Management Lock which is scoped to a Subscription, Resource Group or Resource resource in Azure.



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

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 30
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Static"
    }
    sku = {
      name = "Basic"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "lock" {
  type      = "Microsoft.Authorization/locks@2020-05-01"
  parent_id = azapi_resource.publicIPAddress.id
  name      = var.resource_name
  body = {
    properties = {
      level = "CanNotDelete"
      notes = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Authorization/locks@api-version`. The available api-versions for this resource are: [`2015-01-01`, `2016-09-01`, `2017-04-01`, `2020-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}`  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Authorization/locks?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Authorization/locks/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Authorization/locks/{resourceName}?api-version=2020-05-01
 ```
