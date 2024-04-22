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

resource "azapi_resource" "securityContact" {
  type      = "Microsoft.Security/securityContacts@2017-08-01-preview"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      alertNotifications = "On"
      alertsToAdmins     = "On"
      email              = "basic@example.com"
      phone              = "+1-555-555-5555"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

