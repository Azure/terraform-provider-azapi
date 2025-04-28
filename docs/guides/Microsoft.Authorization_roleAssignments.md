---
subcategory: "Microsoft.Authorization - Azure Resource Manager"
page_title: "roleAssignments"
description: |-
  Manages a Assigns a given Principal (User or Group) to a given Role.
---

# Microsoft.Authorization/roleAssignments - Assigns a given Principal (User or Group) to a given Role

This article demonstrates how to use `azapi` provider to manage the Assigns a given Principal (User or Group) to a given Role resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

data "azurerm_role_definition" "roleAcrpull" {
  name  = "AcrPull"
  scope = azapi_resource.resourceGroup.id
}

resource "azurerm_user_assigned_identity" "uai" {
  name                = "TestUAI"
  resource_group_name = azapi_resource.resourceGroup.name
  location            = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "roleAssignments" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = "6faae21a-0cd6-4536-8c23-a278823d12ed"
  parent_id = azapi_resource.resourceGroup.id
  body = {
    properties = {
      principalId      = azurerm_user_assigned_identity.uai.principal_id
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azurerm_role_definition.roleAcrpull.id
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Authorization/roleAssignments@api-version`. The available api-versions for this resource are: [`2015-07-01`, `2017-10-01-preview`, `2018-01-01-preview`, `2018-09-01-preview`, `2020-03-01-preview`, `2020-04-01-preview`, `2020-08-01-preview`, `2020-10-01-preview`, `2022-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Authorization/roleAssignments?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/roleAssignments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/roleAssignments/{resourceName}?api-version=2022-04-01
 ```
