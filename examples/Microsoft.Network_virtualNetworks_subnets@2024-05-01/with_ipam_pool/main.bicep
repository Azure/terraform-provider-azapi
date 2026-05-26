param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ipamPool 'Microsoft.Network/networkManagers/ipamPools@2024-05-01' = {
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

resource networkManager 'Microsoft.Network/networkManagers@2024-05-01' = {
  location: location
  name: resource_name
  tags: {
    sampleTag: 'sampleTag'
  }
  properties: {
    description: ''
    networkManagerScopeAccesses: []
    networkManagerScopes: {
      managementGroups: []
      subscriptions: [
        '/subscriptions/${data.azapi_client_config.current.subscription_id}'
      ]
    }
  }
}

resource subnet_withIPAM 'Microsoft.Network/virtualNetworks/subnets@2024-05-01' = {
  parent: vnet_withIPAM
  name: resource_name
  properties: {
    ipamPoolPrefixAllocations: [
      {
        numberOfIpAddresses: '100'
        pool: {
          id: ipamPool.id
        }
      }
    ]
  }
}

resource vnet_withIPAM 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      ipamPoolPrefixAllocations: [
        {
          numberOfIpAddresses: '100'
          pool: {
            id: ipamPool.id
          }
        }
      ]
    }
  }
}

