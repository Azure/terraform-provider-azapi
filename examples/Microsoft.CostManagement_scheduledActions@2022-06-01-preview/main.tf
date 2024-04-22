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

resource "azapi_resource" "view" {
  type      = "Microsoft.CostManagement/views@2022-10-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      accumulated = "False"
      chart       = "StackedColumn"
      displayName = "Test View wgvtl"
      kpis = [
        {
          enabled = true
          type    = "Forecast"
        },
      ]
      pivots = [
        {
          name = "ServiceName"
          type = "Dimension"
        },
        {
          name = "ResourceLocation"
          type = "Dimension"
        },
        {
          name = "ResourceGroupName"
          type = "Dimension"
        },
      ]
      query = {
        dataSet = {
          aggregation = {
            totalCost = {
              function = "Sum"
              name     = "Cost"
            }
            totalCostUSD = {
              function = "Sum"
              name     = "CostUSD"
            }
          }
          granularity = "Monthly"
          grouping = [
            {
              name = "ResourceGroupName"
              type = "Dimension"
            },
          ]
          sorting = [
            {
              direction = "Ascending"
              name      = "BillingMonth"
            },
          ]
        }
        timeframe = "MonthToDate"
        type      = "Usage"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "scheduledAction" {
  type      = "Microsoft.CostManagement/scheduledActions@2022-06-01-preview"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    kind = "InsightAlert"
    properties = {
      displayName = "acctest 230630032939736168"
      fileDestination = {
        fileFormats = [
        ]
      }
      notification = {
        message = "Oops, cost anomaly"
        subject = "Hi"
        to = [
          "test@test.com",
          "test@hashicorp.developer",
        ]
      }
      schedule = {
        endDate   = "2024-07-18T11:21:14+00:00"
        frequency = "Daily"
        startDate = "2023-07-18T11:21:14+00:00"
      }
      status = "Enabled"
      viewId = azapi_resource.view.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

