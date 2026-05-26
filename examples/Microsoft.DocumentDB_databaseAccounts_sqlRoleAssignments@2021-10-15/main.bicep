param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.Kusto/clusters@2023-05-02' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    enableAutoStop: true
    enableDiskEncryption: false
    enableDoubleEncryption: false
    enablePurge: false
    enableStreamingIngest: false
    engineType: 'V2'
    publicIPType: 'IPv4'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    trustedExternalTenants: []
  }
  sku: {
    capacity: 1
    name: 'Dev(No SLA)_Standard_D11_v2'
    tier: 'Basic'
  }
}

resource database 'Microsoft.Kusto/clusters/databases@2023-05-02' = {
  parent: cluster
  location: location
  name: resource_name
  kind: 'ReadWrite'
  properties: {}
}

resource databaseAccount 'Microsoft.DocumentDB/databaseAccounts@2021-10-15' = {
  location: location
  name: resource_name
  kind: 'GlobalDocumentDB'
  properties: {
    capabilities: []
    consistencyPolicy: {
      defaultConsistencyLevel: 'Session'
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

resource sqlRoleAssignment 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2021-10-15' = {
  parent: databaseAccount
  name: 'ff419bf7-f8ca-ef51-00d2-3576700c341b'
  properties: {
    principalId: cluster.identity.principalId
    roleDefinitionId: data.azapi_resource.sqlRoleDefinition.id
    scope: databaseAccount.id
  }
}

