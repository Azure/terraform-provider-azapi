param location string = 'westeurope'
param resource_name string = 'acctest0001'
param sql_administrator_password string = null

resource database 'Microsoft.Sql/servers/databases@2021-02-01-preview' = {
  parent: server
  location: location
  name: resource_name
  properties: {
    autoPauseDelay: 0
    collation: 'SQL_Latin1_General_CP1_CI_AS'
    createMode: 'Default'
    elasticPoolId: ''
    highAvailabilityReplicaCount: 0
    isLedgerOn: false
    maintenanceConfigurationId: data.azapi_resource_id.publicMaintenanceConfiguration.id
    minCapacity: 0
    readScale: 'Disabled'
    requestedBackupStorageRedundancy: 'Geo'
    zoneRedundant: false
  }
}

resource jobAgent 'Microsoft.Sql/servers/jobAgents@2020-11-01-preview' = {
  parent: server
  location: location
  name: resource_name
  properties: {
    databaseId: database.id
  }
}

resource server 'Microsoft.Sql/servers@2021-02-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: '4dministr4t0r'
    administratorLoginPassword: sql_administrator_password
    minimalTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    version: '12.0'
  }
}

