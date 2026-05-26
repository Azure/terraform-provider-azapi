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

resource staticMember 'Microsoft.Network/networkManagers/networkGroups/staticMembers@2022-09-01' = {
  parent: networkGroup
  name: resource_name
  properties: {
    resourceId: virtualNetwork.id
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
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

