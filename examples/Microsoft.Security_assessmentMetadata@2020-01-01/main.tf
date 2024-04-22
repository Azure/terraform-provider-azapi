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

resource "azapi_resource" "assessmentMetadatum" {
  type      = "Microsoft.Security/assessmentMetadata@2020-01-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "95c7a001-d595-43af-9754-1310c740d34c"
  body = {
    properties = {
      assessmentType = "CustomerManaged"
      description    = "Test Description"
      displayName    = "Test Display Name"
      severity       = "Medium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

