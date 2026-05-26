param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource agentPool 'Microsoft.ContainerRegistry/registries/agentPools@2019-06-01-preview' = {
  parent: registry
  location: location
  name: resource_name
  properties: {
    count: 1
    os: 'Linux'
    tier: 'S1'
  }
}

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

