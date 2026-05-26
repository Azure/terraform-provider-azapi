param administrator_login string = null
param administrator_login_password string = null
param location string = 'eastus'
param resource_name string = 'acctest0001'

resource flexibleServer 'Microsoft.DBforPostgreSQL/flexibleServers@2023-06-01-preview' = {
  location: location
  name: resource_name
  identity: {
    type: 'None'
    userAssignedIdentities: null
  }
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
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

