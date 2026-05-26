param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource routingIntent 'Microsoft.Network/virtualHubs/routingIntent@2022-09-01' = {
  parent: virtualHub
  name: resource_name
  properties: {
    routingPolicies: [
      {
        destinations: [
          'Internet'
        ]
        name: 'InternetTraffic'
        nextHop: azurerm_firewall.test.id
      }
      {
        destinations: [
          'PrivateTraffic'
        ]
        name: 'PrivateTrafficPolicy'
        nextHop: azurerm_firewall.test.id
      }
    ]
  }
}

resource virtualHub 'Microsoft.Network/virtualHubs@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressPrefix: '10.0.2.0/24'
    hubRoutingPreference: 'ExpressRoute'
    virtualRouterAutoScaleConfiguration: {
      minCapacity: 2
    }
    virtualWan: {
      id: virtualWan.id
    }
  }
}

resource virtualWan 'Microsoft.Network/virtualWans@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    allowBranchToBranchTraffic: true
    disableVpnEncryption: false
    office365LocalBreakoutCategory: 'None'
    type: 'Standard'
  }
}

