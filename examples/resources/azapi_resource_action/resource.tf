terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

variable "enabled" {
  type        = bool
  default     = false
  description = "whether start the spring service"
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "example-spring"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "S0"
}

resource "azapi_resource_action" "start" {
  type                   = "Microsoft.AppPlatform/Spring@2022-05-01-preview"
  resource_id            = azurerm_spring_cloud_service.test.id
  action                 = "start"
  response_export_values = ["*"]

  count = var.enabled ? 1 : 0
}

resource "azapi_resource_action" "stop" {
  type                   = "Microsoft.AppPlatform/Spring@2022-05-01-preview"
  resource_id            = azurerm_spring_cloud_service.test.id
  action                 = "stop"
  response_export_values = ["*"]

  count = var.enabled ? 0 : 1
}
