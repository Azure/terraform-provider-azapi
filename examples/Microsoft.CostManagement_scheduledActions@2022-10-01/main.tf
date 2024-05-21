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

data "azapi_resource_id" "view" {
  type      = "Microsoft.CostManagement/views@2023-04-01-preview"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "ms:CostByService"
}

resource "azapi_resource" "scheduledAction" {
  type      = "Microsoft.CostManagement/scheduledActions@2022-10-01"
  parent_id = data.azapi_resource.subscription.id
  name      = var.resource_name
  body = {
    kind = "Email"
    properties = {
      displayName = "CostByServiceViewerz3k"
      fileDestination = {
        fileFormats = [
        ]
      }
      notification = {
        message = ""
        subject = "Cost Management Report for Terraform Testing on Azure with TTL = 2 Day"
        to = [
          "test@test.com",
          "hashicorp@test.com",
        ]
      }
      notificationEmail = "test@test.com"
      schedule = {
        dayOfMonth   = 0
        daysOfWeek   = null
        endDate      = "2023-07-02T00:00:00Z"
        frequency    = "Daily"
        hourOfDay    = 0
        startDate    = "2023-07-01T00:00:00Z"
        weeksOfMonth = null
      }
      status = "Enabled"
      viewId = data.azapi_resource_id.view.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

