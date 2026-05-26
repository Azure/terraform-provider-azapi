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
    name: 'S1'
  }
}

resource input 'Microsoft.StreamAnalytics/streamingJobs/inputs@2020-03-01' = {
  parent: streamingJob
  name: resource_name
  properties: {
    datasource: {
      properties: {
        consumerGroupName: '$Default'
        endpoint: 'messages/events'
        iotHubNamespace: IotHub.name
        sharedAccessPolicyKey: data.azapi_resource_action.listkeys.output.value[0].primaryKey
        sharedAccessPolicyName: 'iothubowner'
      }
      type: 'Microsoft.Devices/IotHubs'
    }
    serialization: {
      properties: {}
      type: 'Avro'
    }
    type: 'Stream'
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
        query: '   SELECT *\n   INTO [YourOutputAlias]\n   FROM [YourInputAlias]\n'
        streamingUnits: 3
      }
    }
  }
}

