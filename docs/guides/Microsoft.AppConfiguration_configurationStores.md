---
subcategory: "Microsoft.AppConfiguration - Azure App Configuration"
page_title: "configurationStores"
description: |-
  Manages a Azure App Configuration.
---

# Microsoft.AppConfiguration/configurationStores - Azure App Configuration

This article demonstrates how to use `azapi` provider to manage the Azure App Configuration resource in Azure.

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

resource "azapi_resource" "configurationStore" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth      = false
      enablePurgeProtection = false
    }
    sku = {
      name = "standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.AppConfiguration/configurationStores@api-version`. The available api-versions for this resource are: [`2019-02-01-preview`, `2019-10-01`, `2019-11-01-preview`, `2020-06-01`, `2020-07-01-preview`, `2021-03-01-preview`, `2021-10-01-preview`, `2022-03-01-preview`, `2022-05-01`, `2023-03-01`, `2023-08-01-preview`, `2023-09-01-preview`, `2024-05-01`, `2024-06-01`, `2024-06-15-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.AppConfiguration/configurationStores?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppConfiguration/configurationStores/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppConfiguration/configurationStores/{resourceName}?api-version=2024-06-15-preview
 ```
