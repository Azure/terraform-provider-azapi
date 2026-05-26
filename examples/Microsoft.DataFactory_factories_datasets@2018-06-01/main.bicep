param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataset 'Microsoft.DataFactory/factories/datasets@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    description: ''
    linkedServiceName: {
      referenceName: linkedservice.name
      type: 'LinkedServiceReference'
    }
    type: 'Json'
    typeProperties: {
      encodingName: 'UTF-8'
      location: {
        container: 'container'
        fileName: 'bar.txt'
        folderPath: 'foo/bar/'
        type: 'AzureBlobStorageLocation'
      }
    }
  }
}

resource factory 'Microsoft.DataFactory/factories@2018-06-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    repoConfiguration: null
  }
}

resource linkedservice 'Microsoft.DataFactory/factories/linkedservices@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    description: ''
    type: 'AzureBlobStorage'
    typeProperties: {
      serviceEndpoint: storageAccount.properties.primaryEndpoints.blob
    }
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

