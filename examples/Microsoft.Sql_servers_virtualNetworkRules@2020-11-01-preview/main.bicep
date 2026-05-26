param location string = 'westeurope'
param resource_name string = 'acctest0001'
param sql_administrator_password string = null

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'missadmin'
    administratorLoginPassword: sql_administrator_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.7.28.0/25'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: [
      {
        service: 'Microsoft.Sql'
      }
    ]
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.7.28.0/23'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

resource virtualNetworkRule 'Microsoft.Sql/servers/virtualNetworkRules@2020-11-01-preview' = {
  parent: server
  name: resource_name
  properties: {
    ignoreMissingVnetServiceEndpoint: false
    virtualNetworkSubnetId: subnet.id
  }
}

