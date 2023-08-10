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
  body = jsonencode({
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
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "remediation" {
  type      = "Microsoft.PolicyInsights/remediations@2021-10-01"
  parent_id = data.azapi_resource.subscription.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      filters = {
        locations = [
        ]
      }
      policyAssignmentId          = azapi_resource.policyAssignment.id
      policyDefinitionReferenceId = ""
      resourceDiscoveryMode       = "ExistingNonCompliant"
    }
  })
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}

