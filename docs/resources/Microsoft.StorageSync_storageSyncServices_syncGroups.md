---
subcategory: "Microsoft.StorageSync - Storage"
page_title: "storageSyncServices/syncGroups"
description: |-
  Manages a Storage Sync Group.
---

# Microsoft.StorageSync/storageSyncServices/syncGroups - Storage Sync Group

This article demonstrates how to use `azapi` provider to manage the Storage Sync Group resource in Azure.



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

resource "azapi_resource" "storageSyncService" {
  type      = "Microsoft.StorageSync/storageSyncServices@2020-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      incomingTrafficPolicy = "AllowAllTraffic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "syncGroup" {
  type                      = "Microsoft.StorageSync/storageSyncServices/syncGroups@2020-03-01"
  parent_id                 = azapi_resource.storageSyncService.id
  name                      = var.resource_name
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.StorageSync/storageSyncServices/syncGroups@api-version`. The available api-versions for this resource are: [`2017-06-05-preview`, `2018-04-02`, `2018-07-01`, `2018-10-01`, `2019-02-01`, `2019-03-01`, `2019-06-01`, `2019-10-01`, `2020-03-01`, `2020-09-01`, `2022-06-01`, `2022-09-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.StorageSync/storageSyncServices/syncGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{resourceName}/syncGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{resourceName}/syncGroups/{resourceName}?api-version=2022-09-01
 ```
