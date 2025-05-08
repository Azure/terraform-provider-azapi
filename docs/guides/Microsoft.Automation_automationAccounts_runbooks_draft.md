---
subcategory: "Microsoft.Automation - Automation"
page_title: "automationAccounts/runbooks/draft"
description: |-
  Manages a Automation Accounts Runbooks Draft.
---

# Microsoft.Automation/automationAccounts/runbooks/draft - Automation Accounts Runbooks Draft

This article demonstrates how to use `azapi` provider to manage the Automation Accounts Runbooks Draft resource in Azure.

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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
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
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "runbook" {
  type      = "Microsoft.Automation/automationAccounts/runbooks@2019-06-01"
  parent_id = azapi_resource.automationAccount.id
  name      = "Get-AzureVMTutorial"
  location  = var.location
  body = {
    properties = {
      description = "This is a test runbook for terraform acceptance test"
      draft = {
      }
      logActivityTrace = 0
      logProgress      = true
      logVerbose       = true
      runbookType      = "PowerShell"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "draft" {
  type      = "Microsoft.Automation/automationAccounts/runbooks/draft@2018-06-30"
  parent_id = azapi_resource.runbook.id
  name      = "content"
}

resource "azapi_resource_action" "draft" {
  type        = "Microsoft.Automation/automationAccounts/runbooks/draft@2018-06-30"
  resource_id = data.azapi_resource_id.draft.id
  method      = "PUT"
  body = {
    name = "content"
  }
  response_export_values = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Automation/automationAccounts/runbooks/draft@api-version`. The available api-versions for this resource are: [`2015-10-31`, `2018-06-30`, `2019-06-01`, `2022-08-08`, `2023-05-15-preview`, `2023-11-01`, `2024-10-23`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/runbooks/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Automation/automationAccounts/runbooks/draft?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/runbooks/{resourceName}/draft/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/runbooks/{resourceName}/draft/{resourceName}?api-version=2024-10-23
 ```
