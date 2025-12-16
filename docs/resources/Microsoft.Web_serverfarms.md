---
subcategory: "Microsoft.Web - App Service, Azure Functions"
page_title: "serverfarms"
description: |-
  Manages a App Service: Service Plan.
---

# Microsoft.Web/serverfarms - App Service: Service Plan

This article demonstrates how to use `azapi` provider to manage the App Service: Service Plan resource in Azure.



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

resource "azapi_resource" "serverfarm" {
  type      = "Microsoft.Web/serverfarms@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hyperV         = false
      perSiteScaling = false
      reserved       = false
      zoneRedundant  = false
    }
    sku = {
      name = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Web/serverfarms@api-version`. The available api-versions for this resource are: [`2015-08-01`, `2016-09-01`, `2018-02-01`, `2019-08-01`, `2020-06-01`, `2020-09-01`, `2020-10-01`, `2020-12-01`, `2021-01-01`, `2021-01-15`, `2021-02-01`, `2021-03-01`, `2022-03-01`, `2022-09-01`, `2023-01-01`, `2023-12-01`, `2024-04-01`, `2024-11-01`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Web/serverfarms?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/serverfarms/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/serverfarms/{resourceName}?api-version=2025-03-01
 ```
