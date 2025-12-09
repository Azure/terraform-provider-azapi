---
subcategory: "Microsoft.ContainerService - Azure Kubernetes Service (AKS)"
page_title: "fleets"
description: |-
  Manages a Kubernetes Fleet Manager.
---

# Microsoft.ContainerService/fleets - Kubernetes Fleet Manager

This article demonstrates how to use `azapi` provider to manage the Kubernetes Fleet Manager resource in Azure.



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

resource "azapi_resource" "fleet" {
  type      = "Microsoft.ContainerService/fleets@2024-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {}
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerService/fleets@api-version`. The available api-versions for this resource are: [`2022-06-02-preview`, `2022-07-02-preview`, `2022-09-02-preview`, `2023-03-15-preview`, `2023-06-15-preview`, `2023-08-15-preview`, `2023-10-15`, `2024-02-02-preview`, `2024-04-01`, `2024-05-02-preview`, `2025-03-01`, `2025-04-01-preview`, `2025-08-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerService/fleets?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/fleets/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/fleets/{resourceName}?api-version=2025-08-01-preview
 ```
