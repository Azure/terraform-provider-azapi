param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource liveEvent 'Microsoft.Media/mediaServices/liveEvents@2022-08-01' = {
  parent: mediaService
  location: location
  name: resource_name
  properties: {
    input: {
      accessControl: {
        ip: {
          allow: [
            {
              address: '0.0.0.0'
              name: 'AllowAll'
              subnetPrefixLength: 0
            }
          ]
        }
      }
      keyFrameIntervalDuration: 'PT6S'
      streamingProtocol: 'RTMP'
    }
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
  sku: {
    name: 'Standard_GRS'
  }
}

