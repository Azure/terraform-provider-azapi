param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource container 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2023-04-15' = {
  parent: sqlDatabase
  name: resource_name
  properties: {
    options: {}
    resource: {
      id: resource_name
      partitionKey: {
        kind: 'Hash'
        paths: [
          '/definition/id'
        ]
      }
    }
  }
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

resource sqlDatabase 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2021-10-15' = {
  parent: databaseAccount
  name: resource_name
  properties: {
    options: {}
    resource: {
      id: resource_name
    }
  }
}

resource userDefinedFunction 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/userDefinedFunctions@2021-10-15' = {
  parent: container
  name: resource_name
  properties: {
    options: {}
    resource: {
      body: '  \tfunction test() {\n\t\tvar context = getContext();\n\t\tvar response = context.getResponse();\n\t\tresponse.setBody(\'Hello, World\');\n\t}\n'
      id: resource_name
    }
  }
}

