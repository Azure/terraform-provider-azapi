---
subcategory: "Microsoft.Compute - Virtual Machines, Virtual Machine Scale Sets"
page_title: "galleries/applications"
description: |-
  Manages a Gallery Application.
---

# Microsoft.Compute/galleries/applications - Gallery Application

This article demonstrates how to use `azapi` provider to manage the Gallery Application resource in Azure.



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

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "application" {
  type      = "Microsoft.Compute/galleries/applications@2022-03-03"
  parent_id = azapi_resource.gallery.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      supportedOSType = "Linux"
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Compute/galleries/applications@api-version`. The available api-versions for this resource are: [`2019-03-01`, `2019-07-01`, `2019-12-01`, `2020-09-30`, `2021-07-01`, `2021-10-01`, `2022-01-03`, `2022-03-03`, `2022-08-03`, `2023-07-03`, `2024-03-03`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Compute/galleries/applications?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{resourceName}/applications/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{resourceName}/applications/{resourceName}?api-version=2024-03-03
 ```
