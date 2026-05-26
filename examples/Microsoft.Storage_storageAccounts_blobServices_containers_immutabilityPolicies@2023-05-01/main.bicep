param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource container 'Microsoft.Storage/storageAccounts/blobServices/containers@2022-09-01' = {
  name: resource_name
  properties: {
    metadata: {
      key: 'value'
    }
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Standard_LRS'
  }
}

