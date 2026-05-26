param location string = 'westeurope'
param resource_name string = 'acctest0001'
param sql_administrator_password string = null

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'mradministrator'
    administratorLoginPassword: sql_administrator_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

