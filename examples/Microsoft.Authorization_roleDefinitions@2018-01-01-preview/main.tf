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

resource "azapi_resource" "roleDefinition" {
  type      = "Microsoft.Authorization/roleDefinitions@2018-01-01-preview"
  parent_id = data.azapi_resource.subscription.id
  name      = "6faae21a-0cd6-4536-8c23-a278823d12ed"
  body = {
    properties = {
      assignableScopes = [
        data.azapi_resource.subscription.id,
      ]
      description = ""
      permissions = [
        {
          actions = [
            "*",
          ]
          dataActions = [
          ]
          notActions = [
          ]
          notDataActions = [
          ]
        },
      ]
      roleName = var.resource_name
      type     = "CustomRole"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

