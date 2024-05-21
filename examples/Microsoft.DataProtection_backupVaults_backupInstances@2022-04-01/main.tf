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

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforPostgreSQL/servers@2017-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "psqladmin"
      administratorLoginPassword = "H@Sh1CoR3!"
      createMode                 = "Default"
      infrastructureEncryption   = "Disabled"
      minimalTlsVersion          = "TLS1_2"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        backupRetentionDays = 7
        storageAutogrow     = "Enabled"
        storageMB           = 5120
      }
      version = "9.5"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "B_Gen5_2"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "backupVault" {
  type      = "Microsoft.DataProtection/backupVaults@2022-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      storageSettings = [
        {
          datastoreType = "VaultStore"
          type          = "LocallyRedundant"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "backupPolicy" {
  type      = "Microsoft.DataProtection/backupVaults/backupPolicies@2022-04-01"
  parent_id = azapi_resource.backupVault.id
  name      = var.resource_name
  body = {
    properties = {
      datasourceTypes = [
        "Microsoft.DBforPostgreSQL/servers/databases",
      ]
      objectType = "BackupPolicy"
      policyRules = [
        {
          backupParameters = {
            backupType = "Full"
            objectType = "AzureBackupParams"
          }
          dataStore = {
            dataStoreType = "VaultStore"
            objectType    = "DataStoreInfoBase"
          }
          name       = "BackupIntervals"
          objectType = "AzureBackupRule"
          trigger = {
            objectType = "ScheduleBasedTriggerContext"
            schedule = {
              repeatingTimeIntervals = [
                "R/2021-05-23T02:30:00+00:00/P1W",
              ]
            }
            taggingCriteria = [
              {
                isDefault = true
                tagInfo = {
                  id      = "Default_"
                  tagName = "Default"
                }
                taggingPriority = 99
              },
            ]
          }
        },
        {
          isDefault = true
          lifecycles = [
            {
              deleteAfter = {
                duration   = "P4M"
                objectType = "AbsoluteDeleteOption"
              }
              sourceDataStore = {
                dataStoreType = "VaultStore"
                objectType    = "DataStoreInfoBase"
              }
              targetDataStoreCopySettings = [
              ]
            },
          ]
          name       = "Default"
          objectType = "AzureRetentionRule"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.DBforPostgreSQL/servers/databases@2017-12-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      charset   = "UTF8"
      collation = "English_United States.1252"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "backupInstance" {
  type      = "Microsoft.DataProtection/backupVaults/backupInstances@2022-04-01"
  parent_id = azapi_resource.backupVault.id
  name      = var.resource_name
  body = {
    properties = {
      dataSourceInfo = {
        datasourceType   = "Microsoft.DBforPostgreSQL/servers/databases"
        objectType       = "Datasource"
        resourceID       = azapi_resource.database.id
        resourceLocation = azapi_resource.database.location
        resourceName     = azapi_resource.database.name
        resourceType     = "Microsoft.DBforPostgreSQL/servers/databases"
        resourceUri      = ""
      }
      dataSourceSetInfo = {
        datasourceType   = "Microsoft.DBforPostgreSQL/servers/databases"
        objectType       = "DatasourceSet"
        resourceID       = azapi_resource.server.id
        resourceLocation = azapi_resource.server.location
        resourceName     = azapi_resource.server.name
        resourceType     = "Microsoft.DBForPostgreSQL/servers"
        resourceUri      = ""
      }
      datasourceAuthCredentials = null
      friendlyName              = var.resource_name
      objectType                = "BackupInstance"
      policyInfo = {
        policyId = azapi_resource.backupPolicy.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

