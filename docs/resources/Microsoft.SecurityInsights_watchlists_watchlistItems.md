---
subcategory: "Microsoft.SecurityInsights - Microsoft Sentinel"
page_title: "watchlists/watchlistItems"
description: |-
  Manages a Sentinel Watchlist Item.
---

# Microsoft.SecurityInsights/watchlists/watchlistItems - Sentinel Watchlist Item

This article demonstrates how to use `azapi` provider to manage the Sentinel Watchlist Item resource in Azure.

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
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "onboardingState" {
  type      = "Microsoft.SecurityInsights/onboardingStates@2022-11-01"
  parent_id = azapi_resource.workspace.id
  name      = "default"
  body = {
    properties = {
      customerManagedKey = false
    }
  }
}

resource "azapi_resource" "watchlist" {
  type      = "Microsoft.SecurityInsights/watchlists@2022-11-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  body = {
    properties = {
      displayName    = "test"
      itemsSearchKey = "k1"
      provider       = "Microsoft"
      source         = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  depends_on                = [azapi_resource.onboardingState]
}

resource "azapi_resource" "watchlistItem" {
  type      = "Microsoft.SecurityInsights/watchlists/watchlistItems@2022-11-01"
  parent_id = azapi_resource.watchlist.id
  name      = "196abd06-eb4e-4322-9c70-37c32e1a588a"
  body = {
    properties = {
      itemsKeyValue = {
        k1 = "v1"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.SecurityInsights/watchlists/watchlistItems@api-version`. The available api-versions for this resource are: [`2019-01-01-preview`, `2021-03-01-preview`, `2021-04-01`, `2021-09-01-preview`, `2021-10-01`, `2021-10-01-preview`, `2022-01-01-preview`, `2022-04-01-preview`, `2022-05-01-preview`, `2022-06-01-preview`, `2022-07-01-preview`, `2022-08-01`, `2022-08-01-preview`, `2022-09-01-preview`, `2022-10-01-preview`, `2022-11-01`, `2022-11-01-preview`, `2022-12-01-preview`, `2023-02-01`, `2023-02-01-preview`, `2023-03-01-preview`, `2023-04-01-preview`, `2023-05-01-preview`, `2023-06-01-preview`, `2023-07-01-preview`, `2023-08-01-preview`, `2023-09-01-preview`, `2023-10-01-preview`, `2023-11-01`, `2023-12-01-preview`, `2024-01-01-preview`, `2024-03-01`, `2024-04-01-preview`, `2024-09-01`, `2024-10-01-preview`, `2025-01-01-preview`, `2025-03-01`, `2025-04-01-preview`, `2025-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.SecurityInsights/watchlists/watchlistItems?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/watchlistItems/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/watchlistItems/{resourceName}?api-version=2025-06-01
 ```
