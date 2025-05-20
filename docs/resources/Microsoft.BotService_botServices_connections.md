---
subcategory: "Microsoft.BotService - Azure Bot Service"
page_title: "botServices/connections"
description: |-
  Manages a Bot Connection.
---

# Microsoft.BotService/botServices/connections - Bot Connection

This article demonstrates how to use `azapi` provider to manage the Bot Connection resource in Azure.

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
    kind = "bot"
    properties = {
      displayName          = var.resource_name
      isCmekEnabled        = false
      isStreamingSupported = false
      msaAppId             = data.azurerm_client_config.current.tenant_id
    }
    sku = {
      name = "F0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listAuthServiceProviders" {
  type                   = "Microsoft.BotService@2021-05-01-preview"
  resource_id            = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.BotService"
  action                 = "listAuthServiceProviders"
  response_export_values = ["*"]
}

resource "azapi_resource" "connection" {
  type      = "Microsoft.BotService/botServices/connections@2021-05-01-preview"
  parent_id = azapi_resource.botService.id
  name      = var.resource_name
  location  = "global"
  body = {
    kind = "bot"
    properties = {
      clientId          = azapi_resource.botService.output.properties.msaAppId
      clientSecret      = "86546868-e7ed-429f-b0e5-3a1caea7db64"
      scopes            = ""
      serviceProviderId = data.azapi_resource_action.listAuthServiceProviders.output.value[36].properties.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.BotService/botServices/connections@api-version`. The available api-versions for this resource are: [`2017-12-01`, `2018-07-12`, `2020-06-02`, `2021-03-01`, `2021-05-01-preview`, `2022-06-15-preview`, `2022-09-15`, `2023-09-15-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.BotService/botServices/connections?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}/connections/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}/connections/{resourceName}?api-version=2023-09-15-preview
 ```
