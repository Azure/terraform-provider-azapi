---
subcategory: "Microsoft.NotificationHubs - Notification Hubs"
page_title: "namespaces/notificationHubs/authorizationRules"
description: |-
  Manages a Authorization Rule associated with a Notification Hub within a Notification Hub Namespace.
---

# Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules - Authorization Rule associated with a Notification Hub within a Notification Hub Namespace

This article demonstrates how to use `azapi` provider to manage the Authorization Rule associated with a Notification Hub within a Notification Hub Namespace resource in Azure.



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
  type      = "Microsoft.NotificationHubs/namespaces@2017-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enabled       = true
      namespaceType = "NotificationHub"
      region        = "westeurope"
    }
    sku = {
      name = "Free"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "notificationHub" {
  type      = "Microsoft.NotificationHubs/namespaces/notificationHubs@2017-04-01"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules@2017-04-01"
  parent_id = azapi_resource.notificationHub.id
  name      = var.resource_name
  body = {
    properties = {
      rights = [
        "Listen",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules@api-version`. The available api-versions for this resource are: [`2014-09-01`, `2016-03-01`, `2017-04-01`, `2023-01-01-preview`, `2023-09-01`, `2023-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{resourceName}/notificationHubs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{resourceName}/notificationHubs/{resourceName}/authorizationRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{resourceName}/notificationHubs/{resourceName}/authorizationRules/{resourceName}?api-version=2023-10-01-preview
 ```
