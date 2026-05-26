param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ipConfiguration 'Microsoft.Network/virtualHubs/ipConfigurations@2022-07-01' = {
  parent: virtualHub
  name: resource_name
  properties: {
    privateIPAddress: '10.5.1.18'
    privateIPAllocationMethod: 'Static'
    publicIPAddress: {
      id: publicIPAddress.id
    }
    subnet: {
      id: subnet.id
    }
  }
}

resource publicIPAddress 'Microsoft.Network/publicIPAddresses@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    ddosSettings: {
      protectionMode: 'VirtualNetworkInherited'
    }
    idleTimeoutInMinutes: 4
    publicIPAddressVersion: 'IPv4'
    publicIPAllocationMethod: 'Static'
  }
  sku: {
    name: 'Standard'
    tier: 'Regional'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: 'RouteServerSubnet'
  properties: {
    addressPrefix: '10.5.1.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource virtualHub 'Microsoft.Network/virtualHubs@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    hubRoutingPreference: 'ExpressRoute'
    sku: 'Standard'
    virtualRouterAutoScaleConfiguration: {
      minCapacity: 2
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

