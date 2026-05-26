param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource namespace 'Microsoft.ServiceBus/namespaces@2022-01-01-preview' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    publicNetworkAccess: 'Enabled'
    zoneRedundant: false
  }
  sku: {
    capacity: 0
    name: 'Standard'
    tier: 'Standard'
  }
}

resource queue 'Microsoft.ServiceBus/namespaces/queues@2021-06-01-preview' = {
  parent: namespace
  name: resource_name
  properties: {
    deadLetteringOnMessageExpiration: false
    enableBatchedOperations: true
    enableExpress: false
    enablePartitioning: true
    maxDeliveryCount: 10
    maxSizeInMegabytes: 81920
    requiresDuplicateDetection: false
    requiresSession: false
    status: 'Active'
  }
}

