param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource searchService 'Microsoft.Search/searchServices@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    authOptions: {
      apiKeyOnly: {}
    }
    disableLocalAuth: false
    encryptionWithCmk: {
      enforcement: 'Disabled'
    }
    hostingMode: 'default'
    networkRuleSet: {
      ipRules: []
    }
    partitionCount: 1
    publicNetworkAccess: 'Enabled'
    replicaCount: 1
  }
  sku: {
    name: 'standard'
  }
  tags: {
    environment: 'staging'
  }
}

resource sharedPrivateLinkResource 'Microsoft.Search/searchServices/sharedPrivateLinkResources@2022-09-01' = {
  parent: searchService
  name: resource_name
  properties: {
    groupId: 'blob'
    privateLinkResourceId: storageAccount.id
    requestMessage: 'please approve'
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

