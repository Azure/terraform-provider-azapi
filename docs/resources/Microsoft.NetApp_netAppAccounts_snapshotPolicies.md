---
subcategory: "Microsoft.NetApp - Azure NetApp Files"
page_title: "netAppAccounts/snapshotPolicies"
description: |-
  Manages a NetApp Snapshot Policy.
---

# Microsoft.NetApp/netAppAccounts/snapshotPolicies - NetApp Snapshot Policy

This article demonstrates how to use `azapi` provider to manage the NetApp Snapshot Policy resource in Azure.

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
  default = "eastus2"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      activeDirectories = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "snapshotPolicy" {
  type      = "Microsoft.NetApp/netAppAccounts/snapshotPolicies@2022-05-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      dailySchedule = {
        hour            = 22
        minute          = 15
        snapshotsToKeep = 1
      }
      enabled = true
      hourlySchedule = {
        minute          = 15
        snapshotsToKeep = 1
      }
      monthlySchedule = {
        daysOfMonth     = "30,15,1"
        hour            = 5
        minute          = 0
        snapshotsToKeep = 1
      }
      weeklySchedule = {
        day             = "Monday,Friday"
        hour            = 23
        minute          = 0
        snapshotsToKeep = 1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.NetApp/netAppAccounts/snapshotPolicies@api-version`. The available api-versions for this resource are: [`2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-09-01`, `2020-11-01`, `2020-12-01`, `2021-02-01`, `2021-04-01`, `2021-04-01-preview`, `2021-06-01`, `2021-08-01`, `2021-10-01`, `2022-01-01`, `2022-03-01`, `2022-05-01`, `2022-09-01`, `2022-11-01`, `2022-11-01-preview`, `2023-05-01`, `2023-05-01-preview`, `2023-07-01`, `2023-07-01-preview`, `2023-11-01`, `2023-11-01-preview`, `2024-01-01`, `2024-03-01`, `2024-03-01-preview`, `2024-05-01`, `2024-05-01-preview`, `2024-07-01`, `2024-07-01-preview`, `2024-09-01`, `2024-09-01-preview`, `2025-01-01`, `2025-01-01-preview`, `2025-03-01`, `2025-03-01-preview`, `2025-06-01`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.NetApp/netAppAccounts/snapshotPolicies?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/snapshotPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/snapshotPolicies/{resourceName}?api-version=2025-07-01-preview
 ```
