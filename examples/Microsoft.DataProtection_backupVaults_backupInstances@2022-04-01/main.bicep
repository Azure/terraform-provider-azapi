param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource backupInstance 'Microsoft.DataProtection/backupVaults/backupInstances@2022-04-01' = {
  parent: backupVault
  name: resource_name
  properties: {
    dataSourceInfo: {
      datasourceType: 'Microsoft.DBforPostgreSQL/servers/databases'
      objectType: 'Datasource'
      resourceID: database.id
      resourceLocation: database.location
      resourceName: database.name
      resourceType: 'Microsoft.DBforPostgreSQL/servers/databases'
      resourceUri: ''
    }
    dataSourceSetInfo: {
      datasourceType: 'Microsoft.DBforPostgreSQL/servers/databases'
      objectType: 'DatasourceSet'
      resourceID: server.id
      resourceLocation: server.location
      resourceName: server.name
      resourceType: 'Microsoft.DBForPostgreSQL/servers'
      resourceUri: ''
    }
    datasourceAuthCredentials: null
    friendlyName: resource_name
    objectType: 'BackupInstance'
    policyInfo: {
      policyId: backupPolicy.id
    }
  }
}

resource backupPolicy 'Microsoft.DataProtection/backupVaults/backupPolicies@2022-04-01' = {
  parent: backupVault
  name: resource_name
  properties: {
    datasourceTypes: [
      'Microsoft.DBforPostgreSQL/servers/databases'
    ]
    objectType: 'BackupPolicy'
    policyRules: [
      {
        backupParameters: {
          backupType: 'Full'
          objectType: 'AzureBackupParams'
        }
        dataStore: {
          dataStoreType: 'VaultStore'
          objectType: 'DataStoreInfoBase'
        }
        name: 'BackupIntervals'
        objectType: 'AzureBackupRule'
        trigger: {
          objectType: 'ScheduleBasedTriggerContext'
          schedule: {
            repeatingTimeIntervals: [
              'R/2021-05-23T02:30:00+00:00/P1W'
            ]
          }
          taggingCriteria: [
            {
              isDefault: true
              tagInfo: {
                id: 'Default_'
                tagName: 'Default'
              }
              taggingPriority: 99
            }
          ]
        }
      }
      {
        isDefault: true
        lifecycles: [
          {
            deleteAfter: {
              duration: 'P4M'
              objectType: 'AbsoluteDeleteOption'
            }
            sourceDataStore: {
              dataStoreType: 'VaultStore'
              objectType: 'DataStoreInfoBase'
            }
            targetDataStoreCopySettings: []
          }
        ]
        name: 'Default'
        objectType: 'AzureRetentionRule'
      }
    ]
  }
}

resource backupVault 'Microsoft.DataProtection/backupVaults@2022-04-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    storageSettings: [
      {
        datastoreType: 'VaultStore'
        type: 'LocallyRedundant'
      }
    ]
  }
}

resource database 'Microsoft.DBforPostgreSQL/servers/databases@2017-12-01' = {
  parent: server
  name: resource_name
  properties: {
    charset: 'UTF8'
    collation: 'English_United States.1252'
  }
}

resource server 'Microsoft.DBforPostgreSQL/servers@2017-12-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'psqladmin'
    administratorLoginPassword: administrator_login_password
    createMode: 'Default'
    infrastructureEncryption: 'Disabled'
    minimalTlsVersion: 'TLS1_2'
    publicNetworkAccess: 'Enabled'
    sslEnforcement: 'Enabled'
    storageProfile: {
      backupRetentionDays: 7
      storageAutogrow: 'Enabled'
      storageMB: 5120
    }
    version: '9.5'
  }
  sku: {
    capacity: 2
    family: 'Gen5'
    name: 'B_Gen5_2'
    tier: 'Basic'
  }
}

