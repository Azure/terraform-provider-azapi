---
subcategory: "Microsoft.BotService - Azure Bot Service"
page_title: "botServices/channels"
description: |-
  Manages a integration for a Bot Channel.
---

# Microsoft.BotService/botServices/channels - integration for a Bot Channel

This article demonstrates how to use `azapi` provider to manage the integration for a Bot Channel resource in Azure.



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

resource "azapi_resource" "botService" {
  type      = "Microsoft.BotService/botServices@2021-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "bot"
    properties = {
      cmekKeyVaultUrl                   = ""
      description                       = ""
      developerAppInsightKey            = ""
      developerAppInsightsApiKey        = ""
      developerAppInsightsApplicationId = ""
      displayName                       = var.resource_name
      endpoint                          = ""
      iconUrl                           = "https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png"
      isCmekEnabled                     = false
      isStreamingSupported              = false
      msaAppId                          = "12345678-1234-1234-1234-123456789012"
    }
    sku = {
      name = "F0"
    }
  }
  tags = {
    environment = "production"
  }
}

resource "azapi_resource" "channel" {
  type      = "Microsoft.BotService/botServices/channels@2021-05-01-preview"
  parent_id = azapi_resource.botService.id
  name      = "AlexaChannel"
  location  = var.location
  body = {
    kind = "bot"
    properties = {
      channelName = "AlexaChannel"
      properties = {
        alexaSkillId = "amzn1.ask.skill.19126e57-867f-4553-b953-ad0a720dddec"
        isEnabled    = true
      }
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.BotService/botServices/channels@api-version`. The available api-versions for this resource are: [`2017-12-01`, `2018-07-12`, `2020-06-02`, `2021-03-01`, `2021-05-01-preview`, `2022-06-15-preview`, `2022-09-15`, `2023-09-15-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.BotService/botServices/channels?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}/channels/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.BotService/botServices/{resourceName}/channels/{resourceName}?api-version=2023-09-15-preview
 ```
