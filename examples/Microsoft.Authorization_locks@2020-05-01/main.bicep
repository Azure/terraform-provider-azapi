param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource lock 'Microsoft.Authorization/locks@2020-05-01' = {
  parent: publicIPAddress
  name: resource_name
  properties: {
    level: 'CanNotDelete'
    notes: ''
  }
}

resource publicIPAddress 'Microsoft.Network/publicIPAddresses@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    ddosSettings: {
      protectionMode: 'VirtualNetworkInherited'
    }
    idleTimeoutInMinutes: 30
    publicIPAddressVersion: 'IPv4'
    publicIPAllocationMethod: 'Static'
  }
  sku: {
    name: 'Basic'
    tier: 'Regional'
  }
}

