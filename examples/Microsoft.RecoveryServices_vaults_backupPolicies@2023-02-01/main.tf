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

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "backupPolicy" {
  type      = "Microsoft.RecoveryServices/vaults/backupPolicies@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
  body = {
    properties = {
      backupManagementType = "AzureStorage"
      retentionPolicy = {
        dailySchedule = {
          retentionDuration = {
            count        = 10
            durationType = "Days"
          }
          retentionTimes = [
            "2018-07-30T23:00:00Z",
          ]
        }
        retentionPolicyType = "LongTermRetentionPolicy"
      }
      schedulePolicy = {
        schedulePolicyType   = "SimpleSchedulePolicy"
        scheduleRunFrequency = "Daily"
        scheduleRunTimes = [
          "2018-07-30T23:00:00Z",
        ]
      }
      timeZone     = "UTC"
      workLoadType = "AzureFileShare"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

