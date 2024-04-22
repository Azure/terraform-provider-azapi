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

