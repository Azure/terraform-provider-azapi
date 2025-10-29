---
subcategory: "Microsoft.Authorization - Azure Resource Manager"
page_title: "policySetDefinitions"
description: |-
  Manages a policy set definition.
---

# Microsoft.Authorization/policySetDefinitions - policy set definition

This article demonstrates how to use `azapi` provider to manage the policy set definition resource in Azure.



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
  default = "westus"
}

data "azapi_client_config" "current" {}

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

resource "azapi_resource" "policySetDefinition" {
  type      = "Microsoft.Authorization/policySetDefinitions@2025-01-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  name      = "acctestpolset-${var.resource_name}"
  body = {
    properties = {
      description = ""
      displayName = "acctestpolset-${var.resource_name}"
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
      policyDefinitions = [{
        groupNames = []
        parameters = {
          listOfAllowedLocations = {
            value = "[parameters('allowedLocations')]"
          }
        }
        policyDefinitionId          = azapi_resource.policyDefinition.id
        policyDefinitionReferenceId = ""
      }]
      policyType = "Custom"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Authorization/policySetDefinitions@api-version`. The available api-versions for this resource are: [`2017-06-01-preview`, `2018-03-01`, `2018-05-01`, `2019-01-01`, `2019-06-01`, `2019-09-01`, `2020-03-01`, `2020-09-01`, `2021-06-01`, `2023-04-01`, `2024-05-01`, `2025-01-01`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/`  
  `/providers/Microsoft.Management/managementGroups/{managementGroupId}`  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Authorization/policySetDefinitions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Authorization/policySetDefinitions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Authorization/policySetDefinitions/{resourceName}?api-version=2025-03-01
 ```
