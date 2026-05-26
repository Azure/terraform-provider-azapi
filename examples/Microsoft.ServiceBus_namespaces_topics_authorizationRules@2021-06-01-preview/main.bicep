param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.ServiceBus/namespaces/topics/authorizationRules@2021-06-01-preview' = {
  parent: topic
  name: resource_name
  properties: {
    rights: [
      'Send'
    ]
  }
}

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

resource topic 'Microsoft.ServiceBus/namespaces/topics@2021-06-01-preview' = {
  parent: namespace
  name: resource_name
  properties: {
    enableBatchedOperations: false
    enableExpress: false
    enablePartitioning: false
    maxSizeInMegabytes: 5120
    requiresDuplicateDetection: false
    status: 'Active'
    supportOrdering: false
  }
}

