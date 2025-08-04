---
subcategory: "Microsoft.Cache - Azure Cache for Redis"
page_title: "redis/linkedServers"
description: |-
  Manages a Redis Linked Server.
---

# Microsoft.Cache/redis/linkedServers - Redis Linked Server

This article demonstrates how to use `azapi` provider to manage the Redis Linked Server resource in Azure.

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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "resourceGroup_secondary" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "${var.resource_name}-secondary"
  location = var.location
}

resource "azapi_resource" "redis_secondary" {
  type      = "Microsoft.Cache/redis@2024-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-secondary"
  location  = var.location
  body = {
    properties = {
      disableAccessKeyAuthentication = false
      enableNonSslPort               = false
      minimumTlsVersion              = "1.2"
      publicNetworkAccess            = "Enabled"
      redisConfiguration = {
        maxmemory-delta                        = "642"
        maxmemory-policy                       = "allkeys-lru"
        maxmemory-reserved                     = "642"
        preferred-data-persistence-auth-method = ""
      }
      redisVersion = "6"
      sku = {
        capacity = 1
        family   = "P"
        name     = "Premium"
      }
    }
  }
}

resource "azapi_resource" "redis_primary" {
  type      = "Microsoft.Cache/redis@2024-11-01"
  parent_id = azapi_resource.resourceGroup_secondary.id
  name      = "${var.resource_name}-primary"
  location  = var.location
  body = {
    properties = {
      disableAccessKeyAuthentication = false
      enableNonSslPort               = false
      minimumTlsVersion              = "1.2"
      publicNetworkAccess            = "Enabled"
      redisConfiguration = {
        maxmemory-delta                        = "642"
        maxmemory-policy                       = "allkeys-lru"
        maxmemory-reserved                     = "642"
        preferred-data-persistence-auth-method = ""
      }
      redisVersion = "6"
      sku = {
        capacity = 1
        family   = "P"
        name     = "Premium"
      }
    }
  }
}

resource "azapi_resource" "linkedServer" {
  type      = "Microsoft.Cache/redis/linkedServers@2024-11-01"
  parent_id = azapi_resource.redis_primary.id
  name      = "${var.resource_name}-secondary"
  body = {
    properties = {
      linkedRedisCacheId       = azapi_resource.redis_secondary.id
      linkedRedisCacheLocation = var.location
      serverRole               = "Secondary"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Cache/redis/linkedServers@api-version`. The available api-versions for this resource are: [`2017-02-01`, `2017-10-01`, `2018-03-01`, `2019-07-01`, `2020-06-01`, `2020-12-01`, `2021-06-01`, `2022-05-01`, `2022-06-01`, `2023-04-01`, `2023-05-01-preview`, `2023-08-01`, `2024-03-01`, `2024-04-01-preview`, `2024-11-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Cache/redis/linkedServers?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}/linkedServers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}/linkedServers/{resourceName}?api-version=2024-11-01
 ```
