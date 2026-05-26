param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource factory 'Microsoft.DataFactory/factories@2018-06-01' = {
  location: location
  name: resource_name
  properties: {
    globalParameters: {}
    publicNetworkAccess: 'Enabled'
    repoConfiguration: null
  }
}

resource managedVirtualNetwork 'Microsoft.DataFactory/factories/managedVirtualNetworks@2018-06-01' = {
  parent: factory
  name: 'default'
  properties: {}
}

