---
subcategory: "Microsoft.BotService - Azure Bot Service"
page_title: "botServices"
description: |-
  Manages a Azure Bot Service.
---

# Microsoft.BotService/botServices - Azure Bot Service

This article demonstrates how to use `azapi` provider to manage the Azure Bot Service resource in Azure.

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
  default = "westeurope"
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "botService" {
  type      = "Microsoft.BotService/botServices@2021-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    kind = "sdk"
    properties = {
      developerAppInsightKey            = ""
      developerAppInsightsApiKey        = ""
      developerAppInsightsApplicationId = ""
      displayName                       = var.resource_name
      endpoint                          = ""
      luisAppIds = [
      ]
      luisKey  = ""
      msaAppId = data.azurerm_client_config.current.client_id
    }
    sku = {
      name = "F0"
    }
    tags = {
      environment = "production"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.BotService/botServices@api-version`. The available api-versions for this resource are: [`2017-12-01`, `2018-07-12`, `2020-06-02`, `2021-03-01`, `2021-05-01-preview`, `2022-06-15-preview`, `2022-09-15`, `2023-09-15-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.BotService/botServices?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}?api-version=2023-09-15-preview
 ```
