param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource backendAddressPool 'Microsoft.Network/loadBalancers/backendAddressPools@2022-07-01' = {
  parent: loadBalancer
  name: resource_name
  properties: {}
}

resource loadBalancer 'Microsoft.Network/loadBalancers@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    frontendIPConfigurations: [
      {
        name: 'internal'
        properties: {
          publicIPAddress: {
            id: publicIPAddress.id
          }
        }
      }
    ]
  }
  sku: {
    name: 'Standard'
    tier: 'Regional'
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

