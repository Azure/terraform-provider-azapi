param administrator_login string = null
param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource elasticPool 'Microsoft.Sql/servers/elasticPools@2014-04-01' = {
  parent: server
  location: location
  name: resource_name
  properties: {
    dtu: 50
    edition: 'Basic'
    storageMB: 5000
  }
}

resource server 'Microsoft.Sql/servers@2015-05-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
    version: '12.0'
  }
}

