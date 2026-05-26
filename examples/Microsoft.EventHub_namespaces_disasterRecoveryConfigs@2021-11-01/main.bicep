param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource disasterRecoveryConfig 'Microsoft.EventHub/namespaces/disasterRecoveryConfigs@2021-11-01' = {
  parent: namespace
  name: resource_name
  properties: {
    partnerNamespace: namespace2.id
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

resource namespace2 'Microsoft.EventHub/namespaces@2022-01-01-preview' = {
  location: 'westus2'
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

