param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dnsResolver 'Microsoft.Network/dnsResolvers@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    virtualNetwork: {
      id: virtualNetwork.id
    }
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

