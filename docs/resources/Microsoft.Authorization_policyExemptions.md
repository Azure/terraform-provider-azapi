---
subcategory: "Microsoft.Authorization - Azure Resource Manager"
page_title: "policyExemptions"
description: |-
  Manages a Policy Exemption.
---

# Microsoft.Authorization/policyExemptions - Policy Exemption

This article demonstrates how to use `azapi` provider to manage the Policy Exemption resource in Azure.



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
  location  = "westeurope"
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      displayName        = ""
      enforcementMode    = "Default"
      policyDefinitionId = azapi_resource.policyDefinition.id
      scope              = data.azapi_resource.subscription.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "policyExemption" {
  type      = "Microsoft.Authorization/policyExemptions@2020-07-01-preview"
  parent_id = data.azapi_resource.subscription.id
  name      = var.resource_name
  body = {
    properties = {
      exemptionCategory  = "Mitigated"
      policyAssignmentId = azapi_resource.policyAssignment.id
      policyDefinitionReferenceIds = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Authorization/policyExemptions@api-version`. The available api-versions for this resource are: [`2020-07-01-preview`, `2022-07-01-preview`, `2024-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Authorization/policyExemptions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/policyExemptions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Authorization/policyExemptions/{resourceName}?api-version=2024-12-01-preview
 ```
