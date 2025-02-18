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

provider "azurerm" {
  features {
  }
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workflows" {
  type      = "Microsoft.Logic/workflows@2019-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "None"
      userAssignedIdentities = null
    }
    properties = {
      definition = {
        "$schema"      = "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#"
        contentVersion = "1.0.0.0"
      }
      state = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "namespaces" {
  type      = "Microsoft.ServiceBus/namespaces@2022-10-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "None"
      userAssignedIdentities = null
    }
    properties = {
      disableLocalAuth           = false
      minimumTlsVersion          = "1.2"
      premiumMessagingPartitions = 0
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      capacity = 0
      name     = "Basic"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azurerm_managed_api" "test" {
  name     = "servicebus"
  location = var.location

  depends_on = [azapi_resource.workflows, azapi_resource.namespaces]
}

resource "azapi_resource" "connection" {
  type      = "Microsoft.Web/connections@2016-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      api = {
        id = data.azurerm_managed_api.test.id
      }
      displayName = "Service Bus"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
