---
subcategory: "Microsoft.Cache - Azure Cache for Redis"
page_title: "redis"
description: |-
  Manages a Redis Cache.
---

# Microsoft.Cache/redis - Redis Cache

This article demonstrates how to use `azapi` provider to manage the Redis Cache resource in Azure.

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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "redis" {
  type      = "Microsoft.Cache/redis@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        capacity = 2
        family   = "C"
        name     = "Standard"
      }
      enableNonSslPort  = true
      minimumTlsVersion = "1.2"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Cache/redis@api-version`. The available api-versions for this resource are: [`2015-08-01`, `2016-04-01`, `2017-02-01`, `2017-10-01`, `2018-03-01`, `2019-07-01`, `2020-06-01`, `2020-12-01`, `2021-06-01`, `2022-05-01`, `2022-06-01`, `2023-04-01`, `2023-05-01-preview`, `2023-08-01`, `2024-03-01`, `2024-04-01-preview`, `2024-11-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Cache/redis?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}?api-version=2024-11-01
 ```
