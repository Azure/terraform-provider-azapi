param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource consumerGroup 'Microsoft.EventHub/namespaces/eventhubs/consumerGroups@2021-11-01' = {
  parent: eventhub
  name: resource_name
  properties: {
    userMetadata: ''
  }
}

resource eventhub 'Microsoft.EventHub/namespaces/eventhubs@2021-11-01' = {
  parent: namespace
  name: resource_name
  properties: {
    messageRetentionInDays: 1
    partitionCount: 2
    status: 'Active'
  }
}

resource iotConnector 'Microsoft.HealthcareApis/workspaces/iotConnectors@2022-12-01' = {
  parent: workspace
  location: location
  name: resource_name
  properties: {
    deviceMapping: {
      content: {
        template: []
        templateType: 'CollectionContent'
      }
    }
    ingestionEndpointConfiguration: {
      consumerGroup: consumerGroup.id
      eventHubName: eventhub.name
      fullyQualifiedEventHubNamespace: '${namespace.name}.servicebus.windows.net'
    }
  }
}

resource namespace 'Microsoft.EventHub/namespaces@2022-01-01-preview' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    isAutoInflateEnabled: false
    publicNetworkAccess: 'Enabled'
    zoneRedundant: false
  }
  sku: {
    capacity: 1
    name: 'Standard'
    tier: 'Standard'
  }
}

resource workspace 'Microsoft.HealthcareApis/workspaces@2022-12-01' = {
  location: location
  name: resource_name
}

