param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource application 'Microsoft.Batch/batchAccounts/applications@2022-10-01' = {
  parent: batchAccount
  name: resource_name
  properties: {
    allowUpdates: true
    defaultVersion: ''
    displayName: ''
  }
}

resource batchAccount 'Microsoft.Batch/batchAccounts@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    autoStorage: {
      authenticationMode: 'StorageKeys'
      storageAccountId: storageAccount.id
    }
    encryption: {
      keySource: 'Microsoft.Batch'
    }
    poolAllocationMode: 'BatchService'
    publicNetworkAccess: 'Enabled'
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: true
    allowCrossTenantReplication: true
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    encryption: {
      keySource: 'Microsoft.Storage'
      services: {
        queue: {
          keyType: 'Service'
        }
        table: {
          keyType: 'Service'
        }
      }
    }
    isHnsEnabled: false
    isNfsV3Enabled: false
    isSftpEnabled: false
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      defaultAction: 'Allow'
    }
    publicNetworkAccess: 'Enabled'
    supportsHttpsTrafficOnly: true
  }
  sku: {
    name: 'Standard_LRS'
  }
}

