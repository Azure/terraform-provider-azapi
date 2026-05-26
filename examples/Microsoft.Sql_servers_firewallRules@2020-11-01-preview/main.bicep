param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource firewallRule 'Microsoft.Sql/servers/firewallRules@2020-11-01-preview' = {
  parent: server
  name: resource_name
  properties: {
    endIpAddress: '255.255.255.255'
    startIpAddress: '0.0.0.0'
  }
}

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'msincredible'
    administratorLoginPassword: administrator_login_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

