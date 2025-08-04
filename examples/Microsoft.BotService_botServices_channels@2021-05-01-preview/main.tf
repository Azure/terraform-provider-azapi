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

