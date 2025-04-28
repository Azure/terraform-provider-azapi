---
subcategory: "Microsoft.CostManagement - Cost Management"
page_title: "views"
description: |-
  Manages a Azure Cost Management View.
---

# Microsoft.CostManagement/views - Azure Cost Management View

This article demonstrates how to use `azapi` provider to manage the Azure Cost Management View resource in Azure.

## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.CostManagement/views@api-version`. The available api-versions for this resource are: [`2019-04-01-preview`, `2019-11-01`, `2020-06-01`, `2021-10-01`, `2022-08-01-preview`, `2022-10-01`, `2022-10-01-preview`, `2022-10-05-preview`, `2023-03-01`, `2023-04-01-preview`, `2023-07-01-preview`, `2023-08-01`, `2023-09-01`, `2023-11-01`, `2024-08-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.CostManagement/views?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.CostManagement/views/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.CostManagement/views/{resourceName}?api-version=2024-08-01
 ```
