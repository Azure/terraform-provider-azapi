param deploy_locations null = [
  'westeurope'
]
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource networkGroup 'Microsoft.Network/networkManagers/networkGroups@2024-05-01' = {
  parent: networkManager
  name: resource_name
  properties: {}
}

resource networkManager 'Microsoft.Network/networkManagers@2024-05-01' = {
  location: location
  name: resource_name
  properties: {
    description: ''
    networkManagerScopeAccesses: [
      'Routing'
    ]
    networkManagerScopes: {
      managementGroups: []
      subscriptions: [
        '/subscriptions/${data.azapi_client_config.current.subscription_id}'
      ]
    }
  }
}

resource routingConfiguration 'Microsoft.Network/networkManagers/routingConfigurations@2024-05-01' = {
  parent: networkManager
  name: resource_name
  properties: {
    description: 'example routing configuration'
  }
}

resource rule 'Microsoft.Network/networkManagers/routingConfigurations/ruleCollections/rules@2024-05-01' = {
  parent: ruleCollection
  name: resource_name
  properties: {
    description: 'example rule'
    destination: {
      destinationAddress: '10.0.0.0/16'
      type: 'AddressPrefix'
    }
    nextHop: {
      nextHopAddress: ''
      nextHopType: 'VirtualNetworkGateway'
    }
  }
}

resource ruleCollection 'Microsoft.Network/networkManagers/routingConfigurations/ruleCollections@2024-05-01' = {
  parent: routingConfiguration
  name: resource_name
  properties: {
    appliesTo: [
      {
        networkGroupId: networkGroup.id
      }
    ]
    description: 'example rule collection'
  }
}

resource staticMember 'Microsoft.Network/networkManagers/networkGroups/staticMembers@2024-05-01' = {
  parent: networkGroup
  name: resource_name
  properties: {
    resourceId: virtualNetwork.id
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/22'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

