param administrator_login string = null
param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource database 'Microsoft.DBforMariaDB/servers/databases@2018-06-01' = {
  parent: server
  name: resource_name
  properties: {
    charset: 'utf8'
    collation: 'utf8_general_ci'
  }
}

resource server 'Microsoft.DBforMariaDB/servers@2018-06-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
    createMode: 'Default'
    minimalTlsVersion: 'TLS1_2'
    publicNetworkAccess: 'Enabled'
    sslEnforcement: 'Enabled'
    storageProfile: {
      backupRetentionDays: 7
      storageAutogrow: 'Enabled'
      storageMB: 51200
    }
    version: '10.2'
  }
  sku: {
    capacity: 2
    family: 'Gen5'
    name: 'B_Gen5_2'
    tier: 'Basic'
  }
}

