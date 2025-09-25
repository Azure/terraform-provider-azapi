---
subcategory: "Microsoft.StorageMover - Azure Storage Mover"
page_title: "storageMovers/endpoints"
description: |-
  Manages a Storage Mover Source Endpoint.
---

# Microsoft.StorageMover/storageMovers/endpoints - Storage Mover Source Endpoint

This article demonstrates how to use `azapi` provider to manage the Storage Mover Source Endpoint resource in Azure.

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

resource "azapi_resource" "storageMover" {
  type      = "Microsoft.StorageMover/storageMovers@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "endpoint" {
  type      = "Microsoft.StorageMover/storageMovers/endpoints@2023-03-01"
  parent_id = azapi_resource.storageMover.id
  name      = var.resource_name
  body = {
    properties = {
      endpointType = "NfsMount"
      export       = ""
      host         = "192.168.0.1"
      nfsVersion   = "NFSauto"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.StorageMover/storageMovers/endpoints@api-version`. The available api-versions for this resource are: [`2022-07-01-preview`, `2023-03-01`, `2023-07-01-preview`, `2023-10-01`, `2024-07-01`, `2025-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageMover/storageMovers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.StorageMover/storageMovers/endpoints?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageMover/storageMovers/{resourceName}/endpoints/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageMover/storageMovers/{resourceName}/endpoints/{resourceName}?api-version=2025-07-01
 ```
