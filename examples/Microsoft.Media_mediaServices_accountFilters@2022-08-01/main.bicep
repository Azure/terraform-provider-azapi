param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource accountFilter 'Microsoft.Media/mediaServices/accountFilters@2022-08-01' = {
  parent: mediaService
  name: 'Filter-1'
  properties: {
    tracks: []
  }
}

resource mediaService 'Microsoft.Media/mediaServices@2021-11-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    storageAccounts: [
      {
        id: storageAccount.id
        type: 'Primary'
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
    name: 'Standard_GRS'
  }
}

