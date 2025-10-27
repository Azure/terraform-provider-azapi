---
subcategory: "Microsoft.CostManagement - Cost Management"
page_title: "scheduledActions"
description: |-
  Manages a Azure Cost Management Scheduled Action.
---

# Microsoft.CostManagement/scheduledActions - Azure Cost Management Scheduled Action

This article demonstrates how to use `azapi` provider to manage the Azure Cost Management Scheduled Action resource in Azure.



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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.CostManagement/scheduledActions@api-version`. The available api-versions for this resource are: [`2022-04-01-preview`, `2022-06-01-preview`, `2022-10-01`, `2023-03-01`, `2023-04-01-preview`, `2023-07-01-preview`, `2023-08-01`, `2023-09-01`, `2023-11-01`, `2024-08-01`, `2024-10-01-preview`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.CostManagement/scheduledActions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.CostManagement/scheduledActions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.CostManagement/scheduledActions/{resourceName}?api-version=2025-03-01
 ```
