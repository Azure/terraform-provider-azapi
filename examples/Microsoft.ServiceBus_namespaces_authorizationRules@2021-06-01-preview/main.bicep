param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.ServiceBus/namespaces/authorizationRules@2021-06-01-preview' = {
  parent: namespace
  name: resource_name
  properties: {
    rights: [
      'Listen'
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

