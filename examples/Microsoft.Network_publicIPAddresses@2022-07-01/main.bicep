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

