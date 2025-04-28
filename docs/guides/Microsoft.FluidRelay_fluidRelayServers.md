---
subcategory: "Microsoft.FluidRelay - Fluid Relay"
page_title: "fluidRelayServers"
description: |-
  Manages a Fluid Relay Server.
---

# Microsoft.FluidRelay/fluidRelayServers - Fluid Relay Server

This article demonstrates how to use `azapi` provider to manage the Fluid Relay Server resource in Azure.

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

resource "azapi_resource" "fluidRelayServer" {
  type      = "Microsoft.FluidRelay/fluidRelayServers@2022-05-26"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
    tags = {
      foo = "bar"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.FluidRelay/fluidRelayServers@api-version`. The available api-versions for this resource are: [`2021-03-12-preview`, `2021-06-15-preview`, `2021-08-30-preview`, `2021-09-10-preview`, `2022-02-15`, `2022-04-21`, `2022-05-11`, `2022-05-26`, `2022-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.FluidRelay/fluidRelayServers?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.FluidRelay/fluidRelayServers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.FluidRelay/fluidRelayServers/{resourceName}?api-version=2022-06-01
 ```
