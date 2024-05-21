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
  type      = "Microsoft.ServiceBus/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth    = false
      publicNetworkAccess = "Enabled"
      zoneRedundant       = false
    }
    sku = {
      capacity = 0
      name     = "Standard"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "digitalTwinsInstance" {
  type      = "Microsoft.DigitalTwins/digitalTwinsInstances@2020-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "topic" {
  type      = "Microsoft.ServiceBus/namespaces/topics@2021-06-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      enableBatchedOperations    = false
      enableExpress              = false
      enablePartitioning         = false
      maxSizeInMegabytes         = 5120
      requiresDuplicateDetection = false
      status                     = "Active"
      supportOrdering            = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.ServiceBus/namespaces/topics/authorizationRules@2021-06-01-preview"
  parent_id = azapi_resource.topic.id
  name      = var.resource_name
  body = {
    properties = {
      rights = [
        "Send",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.ServiceBus/namespaces/topics/authorizationRules@2021-06-01-preview"
  resource_id            = azapi_resource.authorizationRule.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "endpoint" {
  type      = "Microsoft.DigitalTwins/digitalTwinsInstances/endpoints@2020-12-01"
  parent_id = azapi_resource.digitalTwinsInstance.id
  name      = var.resource_name
  body = {
    properties = {
      authenticationType        = "KeyBased"
      deadLetterSecret          = ""
      endpointType              = "ServiceBus"
      primaryConnectionString   = data.azapi_resource_action.listKeys.output.primaryConnectionString
      secondaryConnectionString = data.azapi_resource_action.listKeys.output.secondaryConnectionString
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.primaryConnectionString, body.properties.secondaryConnectionString]
  }
}

