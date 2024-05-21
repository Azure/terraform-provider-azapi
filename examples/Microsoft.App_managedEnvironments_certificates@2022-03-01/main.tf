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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "sharedKeys" {
  type                   = "Microsoft.OperationalInsights/workspaces@2020-08-01"
  resource_id            = azapi_resource.workspace.id
  action                 = "sharedKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "managedEnvironment" {
  type      = "Microsoft.App/managedEnvironments@2022-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azapi_resource.workspace.output.properties.customerId
          sharedKey  = data.azapi_resource_action.sharedKeys.output.primarySharedKey
        }
      }
      vnetConfiguration = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "certificate" {
  type      = "Microsoft.App/managedEnvironments/certificates@2022-03-01"
  parent_id = azapi_resource.managedEnvironment.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      password = "TestAcc"
      value    = "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"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

