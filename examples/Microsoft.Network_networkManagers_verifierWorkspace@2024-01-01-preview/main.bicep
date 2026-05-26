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
        '/subscriptions/${data.azapi_client_config.current.subscription_id}'
      ]
    }
  }
}

resource verifierWorkspace 'Microsoft.Network/networkManagers/verifierWorkspaces@2024-01-01-preview' = {
  parent: networkManager
  location: location
  name: resource_name
  tags: {
    myTag: 'testTag'
  }
  properties: {
    description: 'A sample workspace'
  }
}

