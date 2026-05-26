param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource databaseAccount 'Microsoft.DocumentDB/databaseAccounts@2021-10-15' = {
  location: location
  name: resource_name
  kind: 'GlobalDocumentDB'
  properties: {
    capabilities: []
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

resource sqlRoleDefinition 'Microsoft.DocumentDB/databaseAccounts/sqlRoleDefinitions@2021-10-15' = {
  parent: databaseAccount
  name: 'c3ce1661-d0b9-3476-0a7c-2654ce2f3055'
  properties: {
    assignableScopes: [
      databaseAccount.id
    ]
    permissions: [
      {
        dataActions: [
          'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read'
        ]
      }
    ]
    roleName: resource_name
    type: 'CustomRole'
  }
}

