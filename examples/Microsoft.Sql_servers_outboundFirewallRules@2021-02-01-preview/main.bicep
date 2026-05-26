param admin_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource outboundFirewallRule 'Microsoft.Sql/servers/outboundFirewallRules@2021-02-01-preview' = {
  parent: server
  name: 'sql230630033612934212.database.windows.net'
  properties: {}
}

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'msincredible'
    administratorLoginPassword: admin_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Enabled'
    version: '12.0'
  }
}

