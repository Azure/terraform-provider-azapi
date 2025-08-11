---
subcategory: "Microsoft.Storage - Storage"
page_title: "storageAccounts/tableServices/tables"
description: |-
  Manages a Table within Azure Storage.
---

# Microsoft.Storage/storageAccounts/tableServices/tables - Table within Azure Storage

This article demonstrates how to use `azapi` provider to manage the Table within Azure Storage resource in Azure.

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

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2022-09-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

data "azapi_resource" "tableService" {
  type      = "Microsoft.Storage/storageAccounts/tableServices@2022-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
}

resource "azapi_resource" "table" {
  type      = "Microsoft.Storage/storageAccounts/tableServices/tables@2022-09-01"
  parent_id = data.azapi_resource.tableService.id
  name      = var.resource_name
  body = {
    properties = {
      signedIdentifiers = []
    }
  }
  response_export_values = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Storage/storageAccounts/tableServices/tables@api-version`. The available api-versions for this resource are: [`2019-06-01`, `2020-08-01-preview`, `2021-01-01`, `2021-02-01`, `2021-04-01`, `2021-06-01`, `2021-08-01`, `2021-09-01`, `2022-05-01`, `2022-09-01`, `2023-01-01`, `2023-04-01`, `2023-05-01`, `2024-01-01`, `2025-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/tableServices/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Storage/storageAccounts/tableServices/tables?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/tableServices/{resourceName}/tables/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/tableServices/{resourceName}/tables/{resourceName}?api-version=2025-01-01
 ```
