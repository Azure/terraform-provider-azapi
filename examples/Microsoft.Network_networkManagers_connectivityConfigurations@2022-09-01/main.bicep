param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource connectivityConfiguration 'Microsoft.Network/networkManagers/connectivityConfigurations@2022-09-01' = {
  parent: networkManager
  name: resource_name
  properties: {
    appliesToGroups: [
      {
        groupConnectivity: 'None'
        isGlobal: 'False'
        networkGroupId: networkGroup.id
        useHubGateway: 'False'
      }
    ]
    connectivityTopology: 'HubAndSpoke'
    deleteExistingPeering: 'False'
    hubs: [
      {
        resourceId: virtualNetwork.id
        resourceType: virtualNetwork.output.type
      }
    ]
    isGlobal: 'False'
  }
}

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
      'Connectivity'
    ]
    networkManagerScopes: {
      managementGroups: []
      subscriptions: [
        data.azapi_resource.subscription.id
      ]
    }
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    flowTimeoutInMinutes: 10
    subnets: []
  }
}

