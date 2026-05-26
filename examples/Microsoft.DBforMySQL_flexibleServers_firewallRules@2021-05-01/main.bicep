param location string = 'westeurope'
param mysql_administrator_password string = null
param resource_name string = 'acctest0001'

resource firewallRule 'Microsoft.DBforMySQL/flexibleServers/firewallRules@2021-05-01' = {
  parent: flexibleServer
  name: resource_name
  properties: {
    endIpAddress: '255.255.255.255'
    startIpAddress: '0.0.0.0'
  }
}

resource flexibleServer 'Microsoft.DBforMySQL/flexibleServers@2021-05-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'adminTerraform'
    administratorLoginPassword: mysql_administrator_password
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
    version: '5.7'
  }
  sku: {
    name: 'Standard_B1s'
    tier: 'Burstable'
  }
}

