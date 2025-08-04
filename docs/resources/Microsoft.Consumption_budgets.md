---
subcategory: "Microsoft.Consumption - Cost Management"
page_title: "budgets"
description: |-
  Manages a Consumption Budget.
---

# Microsoft.Consumption/budgets - Consumption Budget

This article demonstrates how to use `azapi` provider to manage the Consumption Budget resource in Azure.

## Example Usage

### default

```hcl
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

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Consumption/budgets@api-version`. The available api-versions for this resource are: [`2017-12-30-preview`, `2018-01-31`, `2018-03-31`, `2018-06-30`, `2018-08-31`, `2018-10-01`, `2019-01-01`, `2019-04-01-preview`, `2019-05-01`, `2019-05-01-preview`, `2019-06-01`, `2019-10-01`, `2019-11-01`, `2021-05-01`, `2021-10-01`, `2022-09-01`, `2023-03-01`, `2023-05-01`, `2023-11-01`, `2024-08-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Consumption/budgets?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Consumption/budgets/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Consumption/budgets/{resourceName}?api-version=2024-08-01
 ```
