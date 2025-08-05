---
subcategory: "Microsoft.AADIAM - Microsoft Entra ID"
page_title: "diagnosticSettings"
description: |-
  Manages a Azure Active Directory Diagnostic Setting for Azure Monitor.
---

# Microsoft.AADIAM/diagnosticSettings - Azure Active Directory Diagnostic Setting for Azure Monitor

This article demonstrates how to use `azapi` provider to manage the Azure Active Directory Diagnostic Setting for Azure Monitor resource in Azure.

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

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westus"
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-EHN-unique"
  location  = var.location
  body = {
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      minimumTlsVersion    = "1.2"
      publicNetworkAccess  = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Basic"
      tier     = "Basic"
    }
  }
}

resource "azapi_resource" "eventhub" {
  type      = "Microsoft.EventHub/namespaces/eventhubs@2024-01-01"
  parent_id = azapi_resource.namespace.id
  name      = "${var.resource_name}-EH-unique"
  body = {
    properties = {
      messageRetentionInDays = 1
      partitionCount         = 2
      status                 = "Active"
    }
  }
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.EventHub/namespaces/authorizationRules@2024-01-01"
  parent_id = azapi_resource.namespace.id
  name      = "example"
  body = {
    properties = {
      rights = ["Listen", "Send", "Manage"]
    }
  }
}

resource "azapi_resource" "diagnosticSetting" {
  type      = "Microsoft.AADIAM/diagnosticSettings@2017-04-01"
  parent_id = "/"
  name      = "${var.resource_name}-DS-unique"
  body = {
    properties = {
      eventHubAuthorizationRuleId = azapi_resource.authorizationRule.id
      eventHubName                = azapi_resource.eventhub.name
      logs = [
        {
          category = "RiskyUsers"
          enabled  = true
        },
        {
          category = "ServicePrincipalSignInLogs"
          enabled  = true
        },
        {
          category = "SignInLogs"
          enabled  = true
        },
        {
          category = "B2CRequestLogs"
          enabled  = true
        },
        {
          category = "UserRiskEvents"
          enabled  = true
        },
        {
          category = "NonInteractiveUserSignInLogs"
          enabled  = true
        },
        {
          category = "AuditLogs"
          enabled  = true
        }
      ]
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.AADIAM/diagnosticSettings@api-version`. The available api-versions for this resource are: [`2017-04-01`, `2017-04-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.AADIAM/diagnosticSettings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example //providers/Microsoft.AADIAM/diagnosticSettings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example //providers/Microsoft.AADIAM/diagnosticSettings/{resourceName}?api-version=2017-04-01-preview
 ```
