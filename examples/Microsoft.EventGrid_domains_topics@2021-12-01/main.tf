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

resource "azapi_resource" "domain" {
  type      = "Microsoft.EventGrid/domains@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoCreateTopicWithFirstSubscription = true
      autoDeleteTopicWithLastSubscription  = true
      disableLocalAuth                     = false
      inputSchema                          = "EventGridSchema"
      inputSchemaMapping                   = null
      publicNetworkAccess                  = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "topic" {
  type                      = "Microsoft.EventGrid/domains/topics@2021-12-01"
  parent_id                 = azapi_resource.domain.id
  name                      = var.resource_name
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

