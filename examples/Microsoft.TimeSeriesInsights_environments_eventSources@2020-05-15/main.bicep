param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource IotHub 'Microsoft.Devices/IotHubs@2022-04-30-preview' = {
  location: location
  name: resource_name
  properties: {
    cloudToDevice: {}
    enableFileUploadNotifications: false
    messagingEndpoints: {}
    routing: {
      fallbackRoute: {
        condition: 'true'
        endpointNames: [
          'events'
        ]
        isEnabled: true
        source: 'DeviceMessages'
      }
    }
    storageEndpoints: {}
  }
  sku: {
    capacity: 1
    name: 'B1'
  }
  tags: {
    purpose: 'testing'
  }
}

resource environment 'Microsoft.TimeSeriesInsights/environments@2020-05-15' = {
  location: location
  name: resource_name
  kind: 'Gen2'
  properties: {
    storageConfiguration: {
      accountName: storageAccount.name
      managementKey: data.azapi_resource_action.listKeys.output.keys[0].value
    }
    timeSeriesIdProperties: [
      {
        name: 'id'
        type: 'String'
      }
    ]
  }
  sku: {
    capacity: 1
    name: 'L1'
  }
}

resource eventSource 'Microsoft.TimeSeriesInsights/environments/eventSources@2020-05-15' = {
  parent: environment
  location: location
  name: resource_name
  kind: 'Microsoft.IoTHub'
  properties: {
    consumerGroupName: 'test'
    eventSourceResourceId: IotHub.id
    iotHubName: IotHub.name
    keyName: 'iothubowner'
    sharedAccessKey: data.azapi_resource_action.listkeys.output.value[0].primaryKey
    timestampPropertyName: ''
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

