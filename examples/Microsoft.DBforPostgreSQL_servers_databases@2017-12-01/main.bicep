param administrator_login string = null
param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource database 'Microsoft.DBforPostgreSQL/servers/databases@2017-12-01' = {
  parent: server
  name: resource_name
  properties: {
    charset: 'UTF8'
    collation: 'English_United States.1252'
  }
}

resource server 'Microsoft.DBforPostgreSQL/servers@2017-12-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
    createMode: 'Default'
    infrastructureEncryption: 'Disabled'
    minimalTlsVersion: 'TLS1_2'
    publicNetworkAccess: 'Enabled'
    sslEnforcement: 'Enabled'
    storageProfile: {
      backupRetentionDays: 7
      storageAutogrow: 'Enabled'
      storageMB: 51200
    }
    version: '9.6'
  }
  sku: {
    capacity: 2
    family: 'Gen5'
    name: 'GP_Gen5_2'
    tier: 'GeneralPurpose'
  }
}

