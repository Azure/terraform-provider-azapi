param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource account 'Microsoft.Maps/accounts@2021-02-01' = {
  location: 'global'
  name: resource_name
  sku: {
    name: 'G2'
  }
}

resource creator 'Microsoft.Maps/accounts/creators@2021-02-01' = {
  parent: account
  location: location
  name: resource_name
  properties: {
    storageUnits: 1
  }
}

