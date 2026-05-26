param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource bastionHost 'Microsoft.Network/bastionHosts@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    disableCopyPaste: false
    enableFileCopy: false
    enableIpConnect: false
    enableShareableLink: false
    enableTunneling: false
    ipConfigurations: [
      {
        name: 'ip-configuration'
        properties: {
          publicIPAddress: {
            id: publicIPAddress.id
          }
          subnet: {
            id: subnet.id
          }
        }
      }
    ]
    scaleUnits: 2
  }
  sku: {
    name: 'Basic'
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
  name: 'AzureBastionSubnet'
  properties: {
    addressPrefix: '192.168.1.224/27'
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
        '192.168.1.0/24'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

