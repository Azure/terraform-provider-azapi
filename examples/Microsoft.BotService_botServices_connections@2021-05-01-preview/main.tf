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

