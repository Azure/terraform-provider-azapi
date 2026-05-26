param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource networkManager 'Microsoft.Network/networkManagers@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    description: ''
    networkManagerScopeAccesses: [
      'SecurityAdmin'
    ]
    networkManagerScopes: {
      managementGroups: []
      subscriptions: [
        data.azapi_resource.subscription.id
      ]
    }
  }
}

resource scopeConnection 'Microsoft.Network/networkManagers/scopeConnections@2022-09-01' = {
  parent: networkManager
  name: resource_name
  properties: {
    resourceId: data.azapi_resource.subscription.id
    tenantId: data.azapi_resource.subscription.output.tenantId
  }
}

