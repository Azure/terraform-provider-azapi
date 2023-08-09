terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "example-resource-group"
  location = "eastus"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "example"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "test" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azapi_resource" "container_app_environment" {
  type      = "Microsoft.App/managedEnvironments@2022-03-01"
  name      = "example"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azurerm_log_analytics_workspace.test.workspace_id
          sharedKey  = azurerm_log_analytics_workspace.test.primary_shared_key
        }
      }
    }
  })

  // properties/appLogsConfiguration/logAnalyticsConfiguration/sharedKey contains credential which will not be returned,
  // using this property to suppress plan-diff
  ignore_missing_property = true
}

resource "azapi_resource" "dapr" {
  type      = "Microsoft.App/managedEnvironments/daprComponents@2022-03-01"
  name      = "example"
  parent_id = azapi_resource.container_app_environment.id
  body = jsonencode({
    properties = {
      componentType = "state.azure.blobstorage"
      version       = "v1"
      ignoreErrors  = false
      initTimeout   = "5s"
      secrets = [
        {
          name  = "storageaccountkey"
          value = azurerm_storage_account.test.primary_access_key
        }
      ]
      metadata = [
        {
          name  = "accountName"
          value = azurerm_storage_account.test.name
        },
        {
          name      = "accountKey"
          secretRef = "storageaccountkey"
        }
      ]
      scopes = [
        "example"
      ]
    }
  })
  // properties/secrets/value contains credential which will not be returned,
  // using this property to suppress plan-diff
  ignore_missing_property = true
}
