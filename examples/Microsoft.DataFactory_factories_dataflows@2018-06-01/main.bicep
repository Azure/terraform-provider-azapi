param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataflow 'Microsoft.DataFactory/factories/dataflows@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    description: ''
    type: 'Flowlet'
    typeProperties: {
      script: 'source(\n  allowSchemaDrift: true, \n  validateSchema: false, \n  limit: 100, \n  ignoreNoFilesFound: false, \n  documentForm: \'documentPerLine\') ~> source1 \nsource1 sink(\n  allowSchemaDrift: true, \n  validateSchema: false, \n  skipDuplicateMapInputs: true, \n  skipDuplicateMapOutputs: true) ~> sink1\n'
      sinks: [
        {
          description: ''
          linkedService: {
            parameters: {}
            referenceName: linkedservice.name
            type: 'LinkedServiceReference'
          }
          name: 'sink1'
        }
      ]
      sources: [
        {
          description: ''
          linkedService: {
            parameters: {}
            referenceName: linkedservice.name
            type: 'LinkedServiceReference'
          }
          name: 'source1'
        }
      ]
    }
  }
}

resource factory 'Microsoft.DataFactory/factories@2018-06-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    repoConfiguration: null
  }
}

resource linkedservice 'Microsoft.DataFactory/factories/linkedservices@2018-06-01' = {
  parent: factory
  name: resource_name
  properties: {
    description: ''
    type: 'AzureBlobStorage'
    typeProperties: {
      serviceEndpoint: storageAccount.properties.primaryEndpoints.blob
    }
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

