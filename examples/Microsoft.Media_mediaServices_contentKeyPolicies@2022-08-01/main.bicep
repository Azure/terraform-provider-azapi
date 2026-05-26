param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource contentKeyPolicy 'Microsoft.Media/mediaServices/contentKeyPolicies@2022-08-01' = {
  parent: mediaService
  name: resource_name
  properties: {
    description: 'My Policy Description'
    options: [
      {
        configuration: {
          '@odata.type': '#Microsoft.Media.ContentKeyPolicyClearKeyConfiguration'
        }
        name: 'ClearKeyOption'
        restriction: {
          '@odata.type': '#Microsoft.Media.ContentKeyPolicyTokenRestriction'
          audience: 'urn:audience'
          issuer: 'urn:issuer'
          primaryVerificationKey: {
            '@odata.type': '#Microsoft.Media.ContentKeyPolicySymmetricTokenKey'
            keyValue: 'AAAAAAAAAAAAAAAAAAAAAA=='
          }
          requiredClaims: []
          restrictionTokenType: 'Swt'
        }
      }
    ]
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

