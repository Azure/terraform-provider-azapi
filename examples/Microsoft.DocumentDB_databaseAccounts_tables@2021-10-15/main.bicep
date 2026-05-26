param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource databaseAccount 'Microsoft.DocumentDB/databaseAccounts@2021-10-15' = {
  location: location
  name: resource_name
  kind: 'GlobalDocumentDB'
  properties: {
    capabilities: [
      {
        name: 'EnableTable'
      }
    ]
    consistencyPolicy: {
      defaultConsistencyLevel: 'Strong'
      maxIntervalInSeconds: 5
      maxStalenessPrefix: 100
    }
    databaseAccountOfferType: 'Standard'
    defaultIdentity: 'FirstPartyIdentity'
    disableKeyBasedMetadataWriteAccess: false
    disableLocalAuth: false
    enableAnalyticalStorage: false
    enableAutomaticFailover: false
    enableFreeTier: false
    enableMultipleWriteLocations: false
    ipRules: []
    isVirtualNetworkFilterEnabled: false
    locations: [
      {
        failoverPriority: 0
        isZoneRedundant: false
        locationName: 'West Europe'
      }
    ]
    networkAclBypass: 'None'
    networkAclBypassResourceIds: []
    publicNetworkAccess: 'Enabled'
    virtualNetworkRules: []
  }
}

resource table 'Microsoft.DocumentDB/databaseAccounts/tables@2021-10-15' = {
  parent: databaseAccount
  name: resource_name
  properties: {
    options: {}
    resource: {
      id: resource_name
    }
  }
}

