param administrator_login string = null
param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

