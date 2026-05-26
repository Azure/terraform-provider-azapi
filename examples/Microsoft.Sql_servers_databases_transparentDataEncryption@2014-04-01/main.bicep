param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource database 'Microsoft.Sql/servers/databases@2021-02-01-preview' = {
  parent: server
  location: location
  name: resource_name
  properties: {
    autoPauseDelay: 0
    createMode: 'Default'
    elasticPoolId: ''
    highAvailabilityReplicaCount: 0
    isLedgerOn: false
    licenseType: 'LicenseIncluded'
    maintenanceConfigurationId: data.azapi_resource_id.publicMaintenanceConfiguration.id
    minCapacity: 0
    readScale: 'Disabled'
    requestedBackupStorageRedundancy: 'Geo'
    zoneRedundant: false
  }
}

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'mradministrator'
    administratorLoginPassword: administrator_login_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

