param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cach 'Microsoft.StorageCache/caches@2023-01-01' = {
  location: location
  name: resource_name
  properties: {
    cacheSizeGB: 3072
    networkSettings: {
      mtu: 1500
      ntpServer: 'time.windows.com'
    }
    subnet: subnet.id
  }
  sku: {
    name: 'Standard_2G'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.0.2.0/24'
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

