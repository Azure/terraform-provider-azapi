---
subcategory: "Microsoft.Maps - Azure Maps"
page_title: "accounts/creators"
description: |-
  Manages a Azure Maps Creator.
---

# Microsoft.Maps/accounts/creators - Azure Maps Creator

This article demonstrates how to use `azapi` provider to manage the Azure Maps Creator resource in Azure.

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

resource "azapi_resource" "account" {
  type      = "Microsoft.Maps/accounts@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    sku = {
      name = "G2"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "creator" {
  type      = "Microsoft.Maps/accounts/creators@2021-02-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      storageUnits = 1
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Maps/accounts/creators@api-version`. The available api-versions for this resource are: [`2020-02-01-preview`, `2021-02-01`, `2021-07-01-preview`, `2021-12-01-preview`, `2023-06-01`, `2023-08-01-preview`, `2023-12-01-preview`, `2024-01-01-preview`, `2024-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Maps/accounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Maps/accounts/creators?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Maps/accounts/{resourceName}/creators/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Maps/accounts/{resourceName}/creators/{resourceName}?api-version=2024-07-01-preview
 ```
