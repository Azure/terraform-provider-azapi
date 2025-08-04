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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "budget" {
  type      = "Microsoft.Consumption/budgets@2019-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  body = {
    properties = {
      amount   = 1000
      category = "Cost"
      filter = {
        tags = {
          name     = "foo"
          operator = "In"
          values   = ["bar"]
        }
      }
      notifications = {
        "Actual_EqualTo_90.000000_Percent" = {
          contactEmails = ["foo@example.com", "bar@example.com"]
          contactGroups = []
          contactRoles  = []
          enabled       = true
          operator      = "EqualTo"
          threshold     = 90
          thresholdType = "Actual"
        }
      }
      timeGrain = "Monthly"
      timePeriod = {
        startDate = "2025-08-01T00:00:00Z"
      }
    }
  }
}
