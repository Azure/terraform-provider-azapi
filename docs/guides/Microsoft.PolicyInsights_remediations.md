---
subcategory: "Microsoft.PolicyInsights - Azure Policy"
page_title: "remediations"
description: |-
  Manages a Azure Policy Remediation.
---

# Microsoft.PolicyInsights/remediations - Azure Policy Remediation

This article demonstrates how to use `azapi` provider to manage the Azure Policy Remediation resource in Azure.

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

resource "azapi_resource" "policyAssignment" {
  type      = "Microsoft.Authorization/policyAssignments@2022-06-01"
  parent_id = data.azapi_resource.subscription.id
  name      = var.resource_name
  body = {
    properties = {
      displayName     = ""
      enforcementMode = "Default"
      parameters = {
        listOfAllowedLocations = {
          value = [
            "West Europe",
            "West US 2",
            "East US 2",
          ]
        }
      }
      policyDefinitionId = "/providers/Microsoft.Authorization/policyDefinitions/e56962a6-4747-49cd-b67b-bf8b01975c4c"
      scope              = data.azapi_resource.subscription.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "remediation" {
  type      = "Microsoft.PolicyInsights/remediations@2021-10-01"
  parent_id = data.azapi_resource.subscription.id
  name      = var.resource_name
  body = {
    properties = {
      filters = {
        locations = [
        ]
      }
      policyAssignmentId          = azapi_resource.policyAssignment.id
      policyDefinitionReferenceId = ""
      resourceDiscoveryMode       = "ExistingNonCompliant"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.PolicyInsights/remediations@api-version`. The available api-versions for this resource are: [`2018-07-01-preview`, `2019-07-01`, `2021-10-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.PolicyInsights/remediations?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.PolicyInsights/remediations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.PolicyInsights/remediations/{resourceName}?api-version=2024-10-01
 ```
