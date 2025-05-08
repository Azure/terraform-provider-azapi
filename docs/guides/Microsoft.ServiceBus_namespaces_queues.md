---
subcategory: "Microsoft.ServiceBus - Service Bus"
page_title: "namespaces/queues"
description: |-
  Manages a ServiceBus Queue.
---

# Microsoft.ServiceBus/namespaces/queues - ServiceBus Queue

This article demonstrates how to use `azapi` provider to manage the ServiceBus Queue resource in Azure.

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

resource "azapi_resource" "queue" {
  type      = "Microsoft.ServiceBus/namespaces/queues@2021-06-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      deadLetteringOnMessageExpiration = false
      enableBatchedOperations          = true
      enableExpress                    = false
      enablePartitioning               = true
      maxDeliveryCount                 = 10
      maxSizeInMegabytes               = 81920
      requiresDuplicateDetection       = false
      requiresSession                  = false
      status                           = "Active"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ServiceBus/namespaces/queues@api-version`. The available api-versions for this resource are: [`2015-08-01`, `2017-04-01`, `2018-01-01-preview`, `2021-01-01-preview`, `2021-06-01-preview`, `2021-11-01`, `2022-01-01-preview`, `2022-10-01-preview`, `2023-01-01-preview`, `2024-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ServiceBus/namespaces/queues?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/queues/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/queues/{resourceName}?api-version=2024-01-01
 ```
