---
subcategory: "Microsoft.ServiceBus - Service Bus"
page_title: "namespaces/topics/subscriptions"
description: |-
  Manages a ServiceBus Subscription.
---

# Microsoft.ServiceBus/namespaces/topics/subscriptions - ServiceBus Subscription

This article demonstrates how to use `azapi` provider to manage the ServiceBus Subscription resource in Azure.

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

resource "azapi_resource" "namespace" {
  type      = "Microsoft.ServiceBus/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth    = false
      publicNetworkAccess = "Enabled"
      zoneRedundant       = false
    }
    sku = {
      capacity = 0
      name     = "Standard"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "topic" {
  type      = "Microsoft.ServiceBus/namespaces/topics@2021-06-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      enableBatchedOperations    = false
      enableExpress              = false
      enablePartitioning         = false
      maxSizeInMegabytes         = 5120
      requiresDuplicateDetection = false
      status                     = "Active"
      supportOrdering            = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subscription" {
  type      = "Microsoft.ServiceBus/namespaces/topics/subscriptions@2021-06-01-preview"
  parent_id = azapi_resource.topic.id
  name      = var.resource_name
  body = {
    properties = {
      clientAffineProperties = {
      }
      deadLetteringOnFilterEvaluationExceptions = true
      deadLetteringOnMessageExpiration          = false
      enableBatchedOperations                   = false
      isClientAffine                            = false
      maxDeliveryCount                          = 10
      requiresSession                           = false
      status                                    = "Active"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ServiceBus/namespaces/topics/subscriptions@api-version`. The available api-versions for this resource are: [`2015-08-01`, `2017-04-01`, `2018-01-01-preview`, `2021-01-01-preview`, `2021-06-01-preview`, `2021-11-01`, `2022-01-01-preview`, `2022-10-01-preview`, `2023-01-01-preview`, `2024-01-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/topics/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ServiceBus/namespaces/topics/subscriptions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/topics/{resourceName}/subscriptions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/topics/{resourceName}/subscriptions/{resourceName}?api-version=2025-05-01-preview
 ```
