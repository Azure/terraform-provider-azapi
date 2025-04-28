---
subcategory: "Microsoft.Authorization - Azure Resource Manager"
page_title: "roleDefinitions"
description: |-
  Manages a custom Role Definition.
---

# Microsoft.Authorization/roleDefinitions - custom Role Definition

This article demonstrates how to use `azapi` provider to manage the custom Role Definition resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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
  default = "eastus"
}

data "azurerm_client_config" "current" {
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2021-01-01"
  resource_id            = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}

resource "azapi_resource" "roleDefinition" {
  type      = "Microsoft.Authorization/roleDefinitions@2018-01-01-preview"
  parent_id = data.azapi_resource.subscription.id
  name      = "6faae21a-0cd6-4536-8c23-a278823d12ed"
  body = {
    properties = {
      assignableScopes = [
        data.azapi_resource.subscription.id,
      ]
      description = ""
      permissions = [
        {
          actions = [
            "*",
          ]
          dataActions = [
          ]
          notActions = [
          ]
          notDataActions = [
          ]
        },
      ]
      roleName = var.resource_name
      type     = "CustomRole"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Authorization/roleDefinitions@api-version`. The available api-versions for this resource are: [`2015-07-01`, `2018-01-01-preview`, `2022-04-01`, `2022-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Authorization/roleDefinitions?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/roleDefinitions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/roleDefinitions/{resourceName}?api-version=2022-05-01-preview
 ```
