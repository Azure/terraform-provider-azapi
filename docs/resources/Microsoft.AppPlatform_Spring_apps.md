---
subcategory: "Microsoft.AppPlatform - Azure Spring Apps"
page_title: "Spring/apps"
description: |-
  Manages a an Azure Spring Cloud Application.
---

# Microsoft.AppPlatform/Spring/apps - an Azure Spring Cloud Application

This article demonstrates how to use `azapi` provider to manage the an Azure Spring Cloud Application resource in Azure.

!> **Note:** Azure Spring Apps Applications (Microsoft.AppPlatform/Spring/apps) is now deprecated and will be retired on 2028-05-31. See https://aka.ms/asaretirement for more information.

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

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2023-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "E0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "app" {
  type      = "Microsoft.AppPlatform/Spring/apps@2023-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      customPersistentDisks = [
      ]
      enableEndToEndTLS = false
      public            = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.AppPlatform/Spring/apps@api-version`. The available api-versions for this resource are: [`2020-07-01`, `2020-11-01-preview`, `2021-06-01-preview`, `2021-09-01-preview`, `2022-01-01-preview`, `2022-03-01-preview`, `2022-04-01`, `2022-05-01-preview`, `2022-09-01-preview`, `2022-11-01-preview`, `2022-12-01`, `2023-01-01-preview`, `2023-03-01-preview`, `2023-05-01-preview`, `2023-07-01-preview`, `2023-09-01-preview`, `2023-11-01-preview`, `2023-12-01`, `2024-01-01-preview`, `2024-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.AppPlatform/Spring/apps?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{resourceName}/apps/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{resourceName}/apps/{resourceName}?api-version=2024-05-01-preview
 ```
