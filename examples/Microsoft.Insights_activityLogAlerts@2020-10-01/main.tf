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

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "actionGroup" {
  type      = "Microsoft.Insights/actionGroups@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      armRoleReceivers = [
      ]
      automationRunbookReceivers = [
      ]
      azureAppPushReceivers = [
      ]
      azureFunctionReceivers = [
      ]
      emailReceivers = [
      ]
      enabled = true
      eventHubReceivers = [
      ]
      groupShortName = "acctestag1"
      itsmReceivers = [
      ]
      logicAppReceivers = [
      ]
      smsReceivers = [
      ]
      voiceReceivers = [
      ]
      webhookReceivers = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "actionGroup2" {
  type      = "Microsoft.Insights/actionGroups@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      armRoleReceivers = [
      ]
      automationRunbookReceivers = [
      ]
      azureAppPushReceivers = [
      ]
      azureFunctionReceivers = [
      ]
      emailReceivers = [
      ]
      enabled = true
      eventHubReceivers = [
      ]
      groupShortName = "acctestag2"
      itsmReceivers = [
      ]
      logicAppReceivers = [
      ]
      smsReceivers = [
      ]
      voiceReceivers = [
      ]
      webhookReceivers = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "activityLogAlert" {
  type      = "Microsoft.Insights/activityLogAlerts@2020-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      actions = {
        actionGroups = [
          {
            actionGroupId = azapi_resource.actionGroup.id
            webhookProperties = {
            }
          },
          {
            actionGroupId = azapi_resource.actionGroup2.id
            webhookProperties = {
              from = "terraform test"
              to   = "microsoft azure"
            }
          },
        ]
      }
      condition = {
        allOf = [
          {
            equals = "ResourceHealth"
            field  = "category"
          },
          {
            anyOf = [
              {
                equals = "Unavailable"
                field  = "properties.currentHealthStatus"
              },
              {
                equals = "Degraded"
                field  = "properties.currentHealthStatus"
              },
            ]
          },
          {
            anyOf = [
              {
                equals = "Unknown"
                field  = "properties.previousHealthStatus"
              },
              {
                equals = "Available"
                field  = "properties.previousHealthStatus"
              },
            ]
          },
          {
            anyOf = [
              {
                equals = "PlatformInitiated"
                field  = "properties.cause"
              },
              {
                equals = "UserInitiated"
                field  = "properties.cause"
              },
            ]
          },
        ]
      }
      description = "This is just a test acceptance."
      enabled     = true
      scopes = [
        azapi_resource.resourceGroup.id,
        azapi_resource.storageAccount.id,
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

