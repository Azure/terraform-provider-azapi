param administrator_login string = null
param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource database 'Microsoft.DBforMySQL/flexibleServers/databases@2021-05-01' = {
  parent: flexibleServer
  name: resource_name
  properties: {
    charset: 'utf8'
    collation: 'utf8_unicode_ci'
  }
}

resource flexibleServer 'Microsoft.DBforMySQL/flexibleServers@2021-05-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
    backup: {
      backupRetentionDays: 7
      geoRedundantBackup: 'Disabled'
    }
    createMode: ''
    dataEncryption: {
      type: 'SystemManaged'
    }
    highAvailability: {
      mode: 'Disabled'
    }
    network: {}
    version: ''
  }
  sku: {
    name: 'Standard_B1s'
    tier: 'Burstable'
  }
}

