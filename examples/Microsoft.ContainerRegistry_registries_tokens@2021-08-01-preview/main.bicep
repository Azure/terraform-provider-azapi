param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource registry 'Microsoft.ContainerRegistry/registries@2021-08-01-preview' = {
  location: location
  name: resource_name
  properties: {
    adminUserEnabled: true
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

resource token 'Microsoft.ContainerRegistry/registries/tokens@2021-08-01-preview' = {
  parent: registry
  name: resource_name
  properties: {
    scopeMapId: data.azapi_resource_id.scopeMap.id
    status: 'enabled'
  }
}

