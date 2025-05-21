---
subcategory: "Microsoft.Insights - Azure Monitor"
page_title: "diagnosticSettings"
description: |-
  Manages a Insights Diagnostic Settings.
---

# Microsoft.Insights/diagnosticSettings - Insights Diagnostic Settings

This article demonstrates how to use `azapi` provider to manage the Insights Diagnostic Settings resource in Azure.

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
  default = "westeurope"
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2021-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      accessPolicies = [
      ]
      createMode                   = "default"
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      publicNetworkAccess  = "Enabled"
      zoneRedundant        = false
    }
    sku = {
      capacity = 1
      name     = "Basic"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.EventHub/namespaces/authorizationRules@2021-11-01"
  parent_id = azapi_resource.namespace.id
  name      = "example"
  body = {
    properties = {
      rights = [
        "Listen",
        "Send",
        "Manage",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "diagnosticSetting" {
  type      = "Microsoft.Insights/diagnosticSettings@2021-05-01-preview"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
  body = {
    properties = {
      eventHubAuthorizationRuleId = azapi_resource.authorizationRule.id
      eventHubName                = azapi_resource.namespace.name
      logs = [
        {
          categoryGroup = "Audit"
          enabled       = true
          retentionPolicy = {
            days    = 0
            enabled = false
          }
        },
      ]
      metrics = [
        {
          category = "AllMetrics"
          enabled  = true
          retentionPolicy = {
            days    = 0
            enabled = false
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Insights/diagnosticSettings@api-version`. The available api-versions for this resource are: [`2015-07-01`, `2016-09-01`, `2017-05-01-preview`, `2020-01-01-preview`, `2021-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Insights/diagnosticSettings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Insights/diagnosticSettings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Insights/diagnosticSettings/{resourceName}?api-version=2021-05-01-preview
 ```
