param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource credential 'Microsoft.DataFactory/factories/credentials@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    annotations: [
      'test'
    ]
    description: 'this is a test'
    type: 'ManagedIdentity'
    typeProperties: {
      resourceId: userAssignedIdentity.id
    }
  }
}

resource factory 'Microsoft.DataFactory/factories@2018-06-01' = {
  identity: [
    {
      identity_ids: [
        userAssignedIdentity.id
      ]
      type: 'UserAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    repoConfiguration: null
  }
}

resource userAssignedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  location: resourceGroup().location
  name: resource_name
}

