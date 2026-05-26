param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource networkGroup 'Microsoft.Network/networkManagers/networkGroups@2022-09-01' = {
  parent: networkManager
  name: resource_name
  properties: {}
}

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

