param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource Spring 'Microsoft.AppPlatform/Spring@2023-05-01-preview' = {
  location: location
  name: resource_name
  properties: {
    zoneRedundant: false
  }
  sku: {
    name: 'S0'
  }
}

resource app 'Microsoft.AppPlatform/Spring/apps@2023-05-01-preview' = {
  parent: Spring
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    customPersistentDisks: []
    enableEndToEndTLS: false
    public: false
  }
}

resource databaseAccount 'Microsoft.DocumentDB/databaseAccounts@2021-10-15' = {
  location: location
  name: resource_name
  kind: 'GlobalDocumentDB'
  properties: {
    capabilities: []
    consistencyPolicy: {
      defaultConsistencyLevel: 'BoundedStaleness'
      maxIntervalInSeconds: 10
      maxStalenessPrefix: 200
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

resource deployment 'Microsoft.AppPlatform/Spring/apps/deployments@2023-05-01-preview' = {
  parent: app
  name: 'deploy-q4uff'
  properties: {
    deploymentSettings: {
      environmentVariables: {}
      resourceRequests: {
        cpu: '1'
        memory: '1Gi'
      }
    }
    source: {
      jvmOptions: ''
      relativePath: '<default>'
      runtimeVersion: 'Java_8'
      type: 'Jar'
    }
  }
  sku: {
    capacity: 1
    name: 'S0'
    tier: 'Standard'
  }
}

resource linker 'Microsoft.ServiceLinker/linkers@2022-05-01' = {
  parent: deployment
  name: resource_name
  properties: {
    authInfo: {
      authType: 'systemAssignedIdentity'
    }
    clientType: 'none'
    targetService: {
      id: sqlDatabase.id
      resourceProperties: null
      type: 'AzureResource'
    }
  }
}

resource sqlDatabase 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2021-10-15' = {
  parent: databaseAccount
  name: resource_name
  properties: {
    options: {
      throughput: 400
    }
    resource: {
      id: resource_name
    }
  }
}

