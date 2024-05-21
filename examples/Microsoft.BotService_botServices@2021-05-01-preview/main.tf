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

