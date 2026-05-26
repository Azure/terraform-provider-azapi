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

resource rule 'Microsoft.Network/networkManagers/securityAdminConfigurations/ruleCollections/rules@2022-09-01' = {
  parent: ruleCollection
  name: resource_name
  kind: 'Custom'
  properties: {
    access: 'Deny'
    destinationPortRanges: []
    destinations: []
    direction: 'Outbound'
    priority: 1
    protocol: 'Tcp'
    sourcePortRanges: []
    sources: []
  }
}

resource ruleCollection 'Microsoft.Network/networkManagers/securityAdminConfigurations/ruleCollections@2022-09-01' = {
  parent: securityAdminConfiguration
  name: resource_name
  properties: {
    appliesToGroups: [
      {
        networkGroupId: networkGroup.id
      }
    ]
  }
}

resource securityAdminConfiguration 'Microsoft.Network/networkManagers/securityAdminConfigurations@2022-09-01' = {
  parent: networkManager
  name: resource_name
  properties: {
    applyOnNetworkIntentPolicyBasedServices: []
  }
}

