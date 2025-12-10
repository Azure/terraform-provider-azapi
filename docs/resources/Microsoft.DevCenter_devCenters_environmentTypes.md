---
subcategory: "Microsoft.DevCenter - Microsoft Dev Box"
page_title: "devCenters/environmentTypes"
description: |-
  Manages a Dev Center Environment Type.
---

# Microsoft.DevCenter/devCenters/environmentTypes - Dev Center Environment Type

This article demonstrates how to use `azapi` provider to manage the Dev Center Environment Type resource in Azure.



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
  type = string
}

variable "location" {
  type = string
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devCenters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {}
  }
}

resource "azapi_resource" "environmentType" {
  type      = "Microsoft.DevCenter/devCenters/environmentTypes@2025-02-01"
  parent_id = azapi_resource.devCenter.id
  name      = var.resource_name
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DevCenter/devCenters/environmentTypes@api-version`. The available api-versions for this resource are: [`2022-08-01-preview`, `2022-09-01-preview`, `2022-10-12-preview`, `2022-11-11-preview`, `2023-01-01-preview`, `2023-04-01`, `2023-08-01-preview`, `2023-10-01-preview`, `2024-02-01`, `2024-05-01-preview`, `2024-06-01-preview`, `2024-07-01-preview`, `2024-08-01-preview`, `2024-10-01-preview`, `2025-02-01`, `2025-04-01-preview`, `2025-07-01-preview`, `2025-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DevCenter/devCenters/environmentTypes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{resourceName}/environmentTypes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{resourceName}/environmentTypes/{resourceName}?api-version=2025-10-01-preview
 ```
