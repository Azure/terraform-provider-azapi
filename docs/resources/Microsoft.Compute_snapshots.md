---
subcategory: "Microsoft.Compute - Virtual Machines, Virtual Machine Scale Sets"
page_title: "snapshots"
description: |-
  Manages a Disk Snapshot.
---

# Microsoft.Compute/snapshots - Disk Snapshot

This article demonstrates how to use `azapi` provider to manage the Disk Snapshot resource in Azure.



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

resource "azapi_resource" "disk" {
  type      = "Microsoft.Compute/disks@2023-04-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}disk"
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption    = "Empty"
        performancePlus = false
      }
      diskSizeGB = 10
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy        = "AllowAll"
      optimizedForFrequentAttach = false
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "snapshot" {
  type      = "Microsoft.Compute/snapshots@2022-03-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}snapshot"
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption = "Copy"
        sourceUri    = azapi_resource.disk.id
      }
      diskSizeGB          = 20
      incremental         = false
      networkAccessPolicy = "AllowAll"
      publicNetworkAccess = "Enabled"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Compute/snapshots@api-version`. The available api-versions for this resource are: [`2016-04-30-preview`, `2017-03-30`, `2018-04-01`, `2018-06-01`, `2018-09-30`, `2019-03-01`, `2019-07-01`, `2019-11-01`, `2020-05-01`, `2020-06-30`, `2020-09-30`, `2020-12-01`, `2021-04-01`, `2021-08-01`, `2021-12-01`, `2022-03-02`, `2022-07-02`, `2023-01-02`, `2023-04-02`, `2023-10-02`, `2024-03-02`, `2025-01-02`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Compute/snapshots?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{resourceName}?api-version=2025-01-02
 ```
