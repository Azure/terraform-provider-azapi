param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource privateLinkScope 'Microsoft.HybridCompute/privateLinkScopes@2022-11-10' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Disabled'
  }
}

