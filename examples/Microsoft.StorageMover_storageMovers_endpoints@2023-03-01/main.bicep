param location string = 'eastus'
param resource_name string = 'acctest0001'

resource endpoint 'Microsoft.StorageMover/storageMovers/endpoints@2023-03-01' = {
  parent: storageMover
  name: resource_name
  properties: {
    endpointType: 'NfsMount'
    export: ''
    host: '192.168.0.1'
    nfsVersion: 'NFSauto'
  }
}

resource storageMover 'Microsoft.StorageMover/storageMovers@2023-03-01' = {
  location: location
  name: resource_name
  properties: {}
}

