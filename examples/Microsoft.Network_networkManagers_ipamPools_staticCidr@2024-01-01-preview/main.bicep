param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ipamPool 'Microsoft.Network/networkManagers/ipamPools@2024-01-01-preview' = {
  parent: networkManager
  location: location
  name: resource_name
  tags: {
    myTag: 'testTag'
  }
  properties: {
    addressPrefixes: [
      '10.0.0.0/24'
    ]
    description: 'Test description.'
    displayName: 'testDisplayName'
    parentPoolName: ''
  }
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
        '/subscriptions/${data.azapi_client_config.current.subscription_id}'
      ]
    }
  }
}

resource staticCidr 'Microsoft.Network/networkManagers/ipamPools/staticCidrs@2024-01-01-preview' = {
  parent: ipamPool
  name: resource_name
  properties: {
    addressPrefixes: [
      '10.0.0.0/25'
    ]
    description: 'test description'
    numberOfIPAddressesToAllocate: ''
  }
}

