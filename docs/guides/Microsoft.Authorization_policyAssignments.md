---
subcategory: "Microsoft.Authorization - Azure Resource Manager"
page_title: "policyAssignments"
description: |-
  Manages a Policy Assignment.
---

# Microsoft.Authorization/policyAssignments - Policy Assignment

This article demonstrates how to use `azapi` provider to manage the Policy Assignment resource in Azure.

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

resource "azapi_resource" "policyDefinition" {
  type      = "Microsoft.Authorization/policyDefinitions@2021-06-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      displayName = "my-policy-definition"
      mode        = "All"
      parameters = {
        allowedLocations = {
          metadata = {
            description = "The list of allowed locations for resources."
            displayName = "Allowed locations"
            strongType  = "location"
          }
          type = "Array"
        }
      }
      policyRule = {
        if = {
          not = {
            field = "location"
            in    = "[parameters('allowedLocations')]"
          }
        }
        then = {
          effect = "audit"
        }
      }
      policyType = "Custom"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
      policyDefinitionId = azapi_resource.policyDefinition.id
      scope              = data.azapi_resource.subscription.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Authorization/policyAssignments@api-version`. The available api-versions for this resource are: [`2015-10-01-preview`, `2015-11-01`, `2016-04-01`, `2016-12-01`, `2017-06-01-preview`, `2018-03-01`, `2018-05-01`, `2019-01-01`, `2019-06-01`, `2019-09-01`, `2020-03-01`, `2020-09-01`, `2021-06-01`, `2022-06-01`, `2023-04-01`, `2024-04-01`, `2024-05-01`, `2025-01-01`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Authorization/policyAssignments?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/policyAssignments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/policyAssignments/{resourceName}?api-version=2025-03-01
 ```
