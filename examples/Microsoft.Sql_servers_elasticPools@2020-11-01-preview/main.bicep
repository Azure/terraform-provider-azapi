param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource elasticPool 'Microsoft.Sql/servers/elasticPools@2020-11-01-preview' = {
  parent: server
  location: location
  name: resource_name
  properties: {
    maintenanceConfigurationId: data.azapi_resource_id.publicMaintenanceConfiguration.id
    maxSizeBytes: 5242880000
    perDatabaseSettings: {
      maxCapacity: 5
      minCapacity: 0
    }
    zoneRedundant: false
  }
  sku: {
    capacity: 50
    family: ''
    name: 'BasicPool'
    tier: 'Basic'
  }
}

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: '4dm1n157r470r'
    administratorLoginPassword: administrator_login_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

