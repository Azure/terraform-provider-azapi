param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource factory 'Microsoft.DataFactory/factories@2018-06-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    repoConfiguration: null
  }
}

resource integrationRuntime 'Microsoft.DataFactory/factories/integrationRuntimes@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    description: ''
    type: 'SelfHosted'
  }
}

