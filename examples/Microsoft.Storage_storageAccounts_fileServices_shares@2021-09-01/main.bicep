param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource share 'Microsoft.Storage/storageAccounts/fileServices/shares@2022-09-01' = {
  name: resource_name
  properties: {
    accessTier: 'Cool'
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

