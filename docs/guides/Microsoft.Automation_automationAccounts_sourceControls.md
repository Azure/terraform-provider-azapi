---
subcategory: "Microsoft.Automation - Automation"
page_title: "automationAccounts/sourceControls"
description: |-
  Manages a Automation Source Control.
---

# Microsoft.Automation/automationAccounts/sourceControls - Automation Source Control

This article demonstrates how to use `azapi` provider to manage the Automation Source Control resource in Azure.

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

variable "pat" {
  type        = string
  sensitive   = true
  description = "GitHub Personal Access Token"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }

  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_resource" "sourceControl" {
  type      = "Microsoft.Automation/automationAccounts/sourceControls@2023-11-01"
  name      = var.resource_name
  parent_id = azapi_resource.automationAccount.id

  body = {
    properties = {
      repoUrl        = "https://github.com/Azure-Samples/acr-build-helloworld-node.git"
      branch         = "master"
      sourceType     = "GitHub"
      folderPath     = "/"
      autoSync       = false
      publishRunbook = false

      securityToken = {
        tokenType   = "PersonalAccessToken"
        accessToken = var.pat
      }
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Automation/automationAccounts/sourceControls@api-version`. The available api-versions for this resource are: [`2017-05-15-preview`, `2019-06-01`, `2020-01-13-preview`, `2022-08-08`, `2023-05-15-preview`, `2023-11-01`, `2024-10-23`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Automation/automationAccounts/sourceControls?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/sourceControls/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/sourceControls/{resourceName}?api-version=2024-10-23
 ```
