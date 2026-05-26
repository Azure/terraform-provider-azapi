param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Standard_LRS'
  }
}

resource table 'Microsoft.Storage/storageAccounts/tableServices/tables@2022-09-01' = {
  name: resource_name
  properties: {
    signedIdentifiers: []
  }
}

