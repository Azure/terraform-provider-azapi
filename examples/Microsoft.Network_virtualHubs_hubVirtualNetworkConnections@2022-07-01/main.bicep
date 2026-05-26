param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource hubVirtualNetworkConnection 'Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2022-07-01' = {
  parent: virtualHub
  name: resource_name
  properties: {
    enableInternetSecurity: false
    remoteVirtualNetwork: {
      id: virtualNetwork.id
    }
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

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.5.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
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

