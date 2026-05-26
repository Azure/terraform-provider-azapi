param location string = 'westus'
param resource_name string = 'acctest0001'

resource cacheRule 'Microsoft.ContainerRegistry/registries/cacheRules@2023-07-01' = {
  parent: registry
  name: '${resource_name}-cache-rule'
  properties: {
    sourceRepository: 'mcr.microsoft.com/hello-world'
    targetRepository: 'target'
  }
}

resource registry 'Microsoft.ContainerRegistry/registries@2023-11-01-preview' = {
  location: location
  name: '${resource_name}registry'
  properties: {
    adminUserEnabled: false
    anonymousPullEnabled: false
    dataEndpointEnabled: false
    networkRuleBypassOptions: 'AzureServices'
    policies: {
      exportPolicy: {
        status: 'enabled'
      }
      quarantinePolicy: {
        status: 'disabled'
      }
      retentionPolicy: {}
      trustPolicy: {}
    }
    publicNetworkAccess: 'Enabled'
    zoneRedundancy: 'Disabled'
  }
  sku: {
    name: 'Basic'
  }
}

