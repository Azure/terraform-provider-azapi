---
subcategory: "Microsoft.Cache - Azure Cache for Redis"
page_title: "redisEnterprise"
description: |-
  Manages a Redis Enterprise Cluster.
---

# Microsoft.Cache/redisEnterprise - Redis Enterprise Cluster

This article demonstrates how to use `azapi` provider to manage the Redis Enterprise Cluster resource in Azure.

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

resource "azapi_resource" "redisEnterprise" {
  type      = "Microsoft.Cache/redisEnterprise@2022-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      minimumTlsVersion = "1.2"
    }
    sku = {
      capacity = 2
      name     = "Enterprise_E100"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Cache/redisEnterprise@api-version`. The available api-versions for this resource are: [`2020-10-01-preview`, `2021-02-01-preview`, `2021-03-01`, `2021-08-01`, `2022-01-01`, `2022-11-01-preview`, `2023-03-01-preview`, `2023-07-01`, `2023-08-01-preview`, `2023-10-01-preview`, `2023-11-01`, `2024-02-01`, `2024-03-01-preview`, `2024-06-01-preview`, `2024-09-01-preview`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Cache/redisEnterprise?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redisEnterprise/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redisEnterprise/{resourceName}?api-version=2024-10-01
 ```
