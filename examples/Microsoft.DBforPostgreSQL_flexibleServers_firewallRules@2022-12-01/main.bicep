param location string = 'westeurope'
param postgresql_administrator_password string = null
param resource_name string = 'acctest0001'

resource firewallRule 'Microsoft.DBforPostgreSQL/flexibleServers/firewallRules@2022-12-01' = {
  parent: flexibleServer
  name: resource_name
  properties: {
    endIpAddress: '122.122.0.0'
    startIpAddress: '122.122.0.0'
  }
}

resource flexibleServer 'Microsoft.DBforPostgreSQL/flexibleServers@2022-12-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'adminTerraform'
    administratorLoginPassword: postgresql_administrator_password
    availabilityZone: '2'
    backup: {
      geoRedundantBackup: 'Disabled'
    }
    highAvailability: {
      mode: 'Disabled'
    }
    network: {}
    storage: {
      storageSizeGB: 32
    }
    version: '12'
  }
  sku: {
    name: 'Standard_D2s_v3'
    tier: 'GeneralPurpose'
  }
}

