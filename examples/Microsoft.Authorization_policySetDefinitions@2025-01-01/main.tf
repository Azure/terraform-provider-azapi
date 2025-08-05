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

