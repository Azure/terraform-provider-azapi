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

resource "azapi_resource" "namespace" {
  type      = "Microsoft.NotificationHubs/namespaces@2017-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enabled       = true
      namespaceType = "NotificationHub"
      region        = "westeurope"
    }
    sku = {
      name = "Free"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "notificationHub" {
  type      = "Microsoft.NotificationHubs/namespaces/notificationHubs@2017-04-01"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules@2017-04-01"
  parent_id = azapi_resource.notificationHub.id
  name      = var.resource_name
  body = {
    properties = {
      rights = [
        "Listen",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

