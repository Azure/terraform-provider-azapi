param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource registry 'Microsoft.ContainerRegistry/registries@2021-08-01-preview' = {
  location: location
  name: resource_name
  properties: {
    adminUserEnabled: false
    anonymousPullEnabled: false
    dataEndpointEnabled: false
    encryption: {
      status: 'disabled'
    }
    networkRuleBypassOptions: 'AzureServices'
    policies: {
      exportPolicy: {
        status: 'enabled'
      }
      quarantinePolicy: {
        status: 'disabled'
      }
      retentionPolicy: {
        status: 'disabled'
      }
      trustPolicy: {
        status: 'disabled'
      }
    }
    publicNetworkAccess: 'Enabled'
    zoneRedundancy: 'Disabled'
  }
  sku: {
    name: 'Premium'
    tier: 'Premium'
  }
}

resource scopeMap 'Microsoft.ContainerRegistry/registries/scopeMaps@2021-08-01-preview' = {
  parent: registry
  name: resource_name
  properties: {
    actions: [
      'repositories/testrepo/content/read'
    ]
    description: ''
  }
}

