param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource output 'Microsoft.StreamAnalytics/streamingJobs/outputs@2021-10-01-preview' = {
  parent: streamingJob
  name: resource_name
  properties: {
    datasource: {
      properties: {
        accountKey: data.azapi_resource_action.listKeys.output.keys[0].value
        accountName: storageAccount.name
        batchSize: 100
        partitionKey: 'foo'
        rowKey: 'bar'
        table: 'foobar'
      }
      type: 'Microsoft.Storage/Table'
    }
    serialization: null
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: true
    allowCrossTenantReplication: true
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    encryption: {
      keySource: 'Microsoft.Storage'
      services: {
        queue: {
          keyType: 'Service'
        }
        table: {
          keyType: 'Service'
        }
      }
    }
    isHnsEnabled: false
    isNfsV3Enabled: false
    isSftpEnabled: false
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      defaultAction: 'Allow'
    }
    publicNetworkAccess: 'Enabled'
    supportsHttpsTrafficOnly: true
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource streamingJob 'Microsoft.StreamAnalytics/streamingJobs@2020-03-01' = {
  location: location
  name: resource_name
  properties: {
    cluster: {}
    compatibilityLevel: '1.0'
    contentStoragePolicy: 'SystemAccount'
    dataLocale: 'en-GB'
    eventsLateArrivalMaxDelayInSeconds: 60
    eventsOutOfOrderMaxDelayInSeconds: 50
    eventsOutOfOrderPolicy: 'Adjust'
    jobType: 'Cloud'
    outputErrorPolicy: 'Drop'
    sku: {
      name: 'Standard'
    }
    transformation: {
      name: 'main'
      properties: {
        query: '    SELECT *\n    INTO [YourOutputAlias]\n    FROM [YourInputAlias]\n'
        streamingUnits: 3
      }
    }
  }
}

