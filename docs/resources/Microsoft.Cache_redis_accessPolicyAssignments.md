---
subcategory: "Microsoft.Cache - Azure Cache for Redis"
page_title: "redis/accessPolicyAssignments"
description: |-
  Manages a Redis Cache Access Policy Assignment.
---

# Microsoft.Cache/redis/accessPolicyAssignments - Redis Cache Access Policy Assignment

This article demonstrates how to use `azapi` provider to manage the Redis Cache Access Policy Assignment resource in Azure.



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
  features {}
}

provider "azapi" {
  skip_provider_registration = false
}

data "azurerm_client_config" "test" {
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

resource "azapi_resource" "accessPolicyAssignment" {
  type      = "Microsoft.Cache/redis/accessPolicyAssignments@2024-03-01"
  name      = var.resource_name
  parent_id = azapi_resource.redis.id
  body = {
    properties = {
      accessPolicyName = "Data Contributor"
      objectId         = data.azurerm_client_config.test.object_id
      objectIdAlias    = "ServicePrincipal"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Cache/redis/accessPolicyAssignments@api-version`. The available api-versions for this resource are: [`2023-05-01-preview`, `2023-08-01`, `2024-03-01`, `2024-04-01-preview`, `2024-11-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Cache/redis/accessPolicyAssignments?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}/accessPolicyAssignments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redis/{resourceName}/accessPolicyAssignments/{resourceName}?api-version=2024-11-01
 ```
