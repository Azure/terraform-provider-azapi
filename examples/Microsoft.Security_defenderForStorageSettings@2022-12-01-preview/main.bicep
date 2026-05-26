param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  properties: {}
  sku: {
    name: 'Standard_LRS'
  }
}

