param location string = 'centralus'
param resource_name string = 'acctest0001'

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
  name: 'GatewaySubnet'
  properties: {
    addressPrefix: '10.6.1.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.6.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
  tags: {
    SkipASMAzSecPack: 'true'
  }
}

resource virtualNetworkGateway 'Microsoft.Network/virtualNetworkGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    activeActive: false
    enableBgp: false
    enablePrivateIpAddress: false
    gatewayType: 'ExpressRoute'
    ipConfigurations: [
      {
        name: 'vnetGatewayConfig'
        properties: {
          privateIPAllocationMethod: 'Dynamic'
          publicIPAddress: {
            id: publicIPAddress.id
          }
          subnet: {
            id: subnet.id
          }
        }
      }
    ]
    sku: {
      name: 'Standard'
      tier: 'Standard'
    }
    vpnType: 'RouteBased'
  }
}

