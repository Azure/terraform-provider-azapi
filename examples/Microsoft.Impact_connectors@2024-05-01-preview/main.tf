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
  default = "westeurope"
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "connector" {
  type      = "Microsoft.Impact/connectors@2024-05-01-preview"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      connectorType = "AzureMonitor"
    }
  }
  schema_validation_enabled = false
}
