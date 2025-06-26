---
subcategory: "Microsoft.DataProtection - Data Protection"
page_title: "backupVaults/backupInstances"
description: |-
  Manages a Backup Instance.
---

# Microsoft.DataProtection/backupVaults/backupInstances - Backup Instance

This article demonstrates how to use `azapi` provider to manage the Backup Instance resource in Azure.

## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DataProtection/backupVaults/backupInstances@api-version`. The available api-versions for this resource are: [`2021-01-01`, `2021-02-01-preview`, `2021-06-01-preview`, `2021-07-01`, `2021-10-01-preview`, `2021-12-01-preview`, `2022-01-01`, `2022-02-01-preview`, `2022-03-01`, `2022-03-31-preview`, `2022-04-01`, `2022-05-01`, `2022-09-01-preview`, `2022-10-01-preview`, `2022-11-01-preview`, `2022-12-01`, `2023-01-01`, `2023-04-01-preview`, `2023-05-01`, `2023-06-01-preview`, `2023-08-01-preview`, `2023-11-01`, `2023-12-01`, `2024-02-01-preview`, `2024-03-01`, `2024-04-01`, `2025-01-01`, `2025-02-01`, `2025-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DataProtection/backupVaults/backupInstances?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{resourceName}/backupInstances/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{resourceName}/backupInstances/{resourceName}?api-version=2025-07-01
 ```
