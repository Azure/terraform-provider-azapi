---
subcategory: "Microsoft.SignalRService - Azure SignalR Service"
page_title: "webPubSub/hubs"
description: |-
  Manages a hub settings for a Web Pubsub service.
---

# Microsoft.SignalRService/webPubSub/hubs - hub settings for a Web Pubsub service

This article demonstrates how to use `azapi` provider to manage the hub settings for a Web Pubsub service resource in Azure.



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

resource "azapi_resource" "webPubSub" {
  type      = "Microsoft.SignalRService/webPubSub@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      disableAadAuth      = false
      disableLocalAuth    = false
      publicNetworkAccess = "Enabled"
      tls = {
        clientCertEnabled = false
      }
    }
    sku = {
      capacity = 1
      name     = "Standard_S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "hub" {
  type      = "Microsoft.SignalRService/webPubSub/hubs@2023-02-01"
  parent_id = azapi_resource.webPubSub.id
  name      = var.resource_name
  body = {
    properties = {
      anonymousConnectPolicy = "Deny"
      eventListeners = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.SignalRService/webPubSub/hubs@api-version`. The available api-versions for this resource are: [`2021-10-01`, `2022-08-01-preview`, `2023-02-01`, `2023-03-01-preview`, `2023-06-01-preview`, `2023-08-01-preview`, `2024-01-01-preview`, `2024-03-01`, `2024-04-01-preview`, `2024-08-01-preview`, `2024-10-01-preview`, `2025-01-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SignalRService/webPubSub/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.SignalRService/webPubSub/hubs?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SignalRService/webPubSub/{resourceName}/hubs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SignalRService/webPubSub/{resourceName}/hubs/{resourceName}?api-version=2025-01-01-preview
 ```
