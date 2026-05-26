param location string = 'westus'
param resource_name string = 'acctest0001'

resource connectedRegistry 'Microsoft.ContainerRegistry/registries/connectedRegistries@2023-11-01-preview' = {
  parent: registry
  name: '${resource_name}connectedregistry'
  properties: {
    clientTokenIds: null
    logging: {
      auditLogStatus: 'Disabled'
      logLevel: 'None'
    }
    mode: 'ReadWrite'
    parent: {
      syncProperties: {
        messageTtl: 'P1D'
        schedule: '* * * * *'
        syncWindow: ''
        tokenId: token.id
      }
    }
  }
}

resource registry 'Microsoft.ContainerRegistry/registries@2023-11-01-preview' = {
  location: location
  name: '${resource_name}registry'
  properties: {
    adminUserEnabled: false
    anonymousPullEnabled: false
    dataEndpointEnabled: true
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
    name: 'Premium'
  }
}

resource scopeMap 'Microsoft.ContainerRegistry/registries/scopeMaps@2023-11-01-preview' = {
  parent: registry
  name: '${resource_name}scopemap'
  properties: {
    actions: [
      'repositories/hello-world/content/delete'
      'repositories/hello-world/content/read'
      'repositories/hello-world/content/write'
      'repositories/hello-world/metadata/read'
      'repositories/hello-world/metadata/write'
      'gateway/${resource_name}connectedregistry/config/read'
      'gateway/${resource_name}connectedregistry/config/write'
      'gateway/${resource_name}connectedregistry/message/read'
      'gateway/${resource_name}connectedregistry/message/write'
    ]
    description: ''
  }
}

resource token 'Microsoft.ContainerRegistry/registries/tokens@2023-11-01-preview' = {
  parent: registry
  name: '${resource_name}token'
  properties: {
    scopeMapId: scopeMap.id
    status: 'enabled'
  }
}

