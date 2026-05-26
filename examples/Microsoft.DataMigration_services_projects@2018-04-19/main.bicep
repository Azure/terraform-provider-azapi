param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource project 'Microsoft.DataMigration/services/projects@2018-04-19' = {
  parent: service
  location: location
  name: resource_name
  properties: {
    sourcePlatform: 'SQL'
    targetPlatform: 'SQLDB'
  }
}

resource service 'Microsoft.DataMigration/services@2018-04-19' = {
  location: location
  name: resource_name
  kind: 'Cloud'
  properties: {
    virtualSubnetId: subnet.id
  }
  sku: {
    name: 'Standard_1vCores'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
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

