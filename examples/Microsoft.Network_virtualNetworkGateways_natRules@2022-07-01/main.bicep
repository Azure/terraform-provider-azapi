param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource natRule 'Microsoft.Network/virtualNetworkGateways/natRules@2022-07-01' = {
  parent: virtualNetworkGateway
  name: resource_name
  properties: {
    externalMappings: [
      {
        addressSpace: '10.1.0.0/26'
      }
    ]
    internalMappings: [
      {
        addressSpace: '10.3.0.0/26'
      }
    ]
    mode: 'EgressSnat'
    type: 'Static'
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
    publicIPAllocationMethod: 'Dynamic'
  }
  sku: {
    name: 'Basic'
    tier: 'Regional'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: 'GatewaySubnet'
  properties: {
    addressPrefix: '10.0.1.0/24'
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
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

resource virtualNetworkGateway 'Microsoft.Network/virtualNetworkGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    activeActive: false
    enableBgp: false
    enablePrivateIpAddress: false
    gatewayType: 'Vpn'
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
      name: 'Basic'
      tier: 'Basic'
    }
    vpnType: 'RouteBased'
  }
}

