param location string = 'eastus'
param resource_name string = 'acctest0001'

resource storageMover 'Microsoft.StorageMover/storageMovers@2023-03-01' = {
  location: location
  name: resource_name
  properties: {}
}

