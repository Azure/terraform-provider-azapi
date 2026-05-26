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

resource pipeline 'Microsoft.DataFactory/factories/pipelines@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    annotations: []
    description: ''
    parameters: {
      test: {
        defaultValue: 'testparameter'
        type: 'String'
      }
    }
    variables: {}
  }
}

