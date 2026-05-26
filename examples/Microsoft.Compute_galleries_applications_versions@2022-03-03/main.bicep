param location string = 'westus'
param resource_name string = 'acctest0001'

resource application 'Microsoft.Compute/galleries/applications@2022-03-03' = {
  parent: gallery
  location: location
  name: '${resource_name}-app'
  properties: {
    supportedOSType: 'Linux'
  }
}

resource container 'Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01' = {
  name: 'mycontainer'
  properties: {
    publicAccess: 'Blob'
  }
}

resource gallery 'Microsoft.Compute/galleries@2022-03-03' = {
  location: location
  name: '${resource_name}sig'
  properties: {
    description: ''
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2023-05-01' = {
  location: location
  name: '${resource_name}acc'
  kind: 'StorageV2'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: true
    allowCrossTenantReplication: false
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    dnsEndpointType: 'Standard'
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
    isLocalUserEnabled: true
    isNfsV3Enabled: false
    isSftpEnabled: false
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      bypass: 'AzureServices'
      defaultAction: 'Allow'
      ipRules: []
      resourceAccessRules: []
      virtualNetworkRules: []
    }
    publicNetworkAccess: 'Enabled'
    supportsHttpsTrafficOnly: true
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource version 'Microsoft.Compute/galleries/applications/versions@2022-03-03' = {
  parent: application
  location: location
  name: '0.0.1'
  properties: {
    publishingProfile: {
      enableHealthCheck: false
      excludeFromLatest: false
      manageActions: {
        install: '[install command]'
        remove: '[remove command]'
        update: ''
      }
      source: {
        defaultConfigurationLink: ''
        mediaLink: 'https://${storageAccount.name}.blob.core.windows.net/mycontainer/myblob'
      }
      targetRegions: [
        {
          name: location
          regionalReplicaCount: 1
          storageAccountType: 'Standard_LRS'
        }
      ]
    }
    safetyProfile: {
      allowDeletionOfReplicatedLocations: true
    }
  }
}

