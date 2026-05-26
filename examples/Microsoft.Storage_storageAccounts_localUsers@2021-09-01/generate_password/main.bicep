param location string = 'westeurope'
param resource_name string = 'acctest0001dfdg'

resource localUser 'Microsoft.Storage/storageAccounts/localUsers@2021-09-01' = {
  parent: storageAccount
  name: resource_name
  properties: {
    hasSharedKey: true
    hasSshKey: false
    hasSshPassword: true
    homeDirectory: 'containername/'
    permissionScopes: [
      {
        permissions: 'cwl'
        resourceName: 'containername'
        service: 'blob'
      }
    ]
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

