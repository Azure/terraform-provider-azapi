---
subcategory: "Microsoft.DataFactory - Data Factory"
page_title: "factories/credentials"
description: |-
  Manages a Data Factory credential resource.
---

# Microsoft.DataFactory/factories/credentials - Data Factory credential resource

This article demonstrates how to use `azapi` provider to manage the Data Factory credential resource resource in Azure.

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

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "credential" {
  type      = "Microsoft.DataFactory/factories/credentials@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      type        = "ManagedIdentity"
      annotations = ["test"]
      description = "this is a test"
      typeProperties = {
        resourceId = azapi_resource.userAssignedIdentity.id
      }
    }

  }
  ignore_casing           = false
  ignore_missing_property = false
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DataFactory/factories/credentials@api-version`. The available api-versions for this resource are: [`2018-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DataFactory/factories/credentials?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}/credentials/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}/credentials/{resourceName}?api-version=2018-06-01
 ```
