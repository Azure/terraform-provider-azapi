param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.EventHub/namespaces/authorizationRules@2021-11-01' = {
  parent: namespace
  name: resource_name
  properties: {
    rights: [
      'Listen'
      'Send'
      'Manage'
    ]
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
    name: 'Basic'
    tier: 'Basic'
  }
}

