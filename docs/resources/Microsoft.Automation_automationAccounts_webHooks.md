---
subcategory: "Microsoft.Automation - Automation"
page_title: "automationAccounts/webHooks"
description: |-
  Manages a Automation Runbook's Webhook.
---

# Microsoft.Automation/automationAccounts/webHooks - Automation Runbook's Webhook

This article demonstrates how to use `azapi` provider to manage the Automation Runbook's Webhook resource in Azure.



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

resource "azapi_resource_action" "publish_runbook" {
  type        = "Microsoft.Automation/automationAccounts/runbooks@2022-08-08"
  resource_id = azapi_resource.runbook.id
  action      = "publish"
}

resource "azapi_resource" "webHook" {
  type      = "Microsoft.Automation/automationAccounts/webHooks@2015-10-31"
  parent_id = azapi_resource.automationAccount.id
  name      = "TestRunbook_webhook"
  body = {
    properties = {
      expiryTime = "2025-06-30T04:27:24Z"
      isEnabled  = true
      parameters = {
      }
      runOn = ""
      runbook = {
        name = azapi_resource.runbook.name
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  depends_on                = [azapi_resource_action.publish_runbook]
  lifecycle {
    ignore_changes = [body.properties.expiryTime]
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Automation/automationAccounts/webHooks@api-version`. The available api-versions for this resource are: [`2015-10-31`, `2023-05-15-preview`, `2024-10-23`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Automation/automationAccounts/webHooks?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/webHooks/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/webHooks/{resourceName}?api-version=2024-10-23
 ```
