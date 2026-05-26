param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource Spring 'Microsoft.AppPlatform/Spring@2023-05-01-preview' = {
  location: location
  name: resource_name
  properties: {
    zoneRedundant: false
  }
  sku: {
    name: 'S0'
  }
}

resource app 'Microsoft.AppPlatform/Spring/apps@2023-05-01-preview' = {
  parent: Spring
  location: location
  name: resource_name
  properties: {
    customPersistentDisks: []
    enableEndToEndTLS: false
    public: false
  }
}

resource binding 'Microsoft.AppPlatform/Spring/apps/bindings@2023-05-01-preview' = {
  parent: app
  name: resource_name
  properties: {
    bindingParameters: {
      useSsl: 'true'
    }
    key: data.azapi_resource_action.listKeys.output.primaryKey
    resourceId: redis.id
  }
}

resource redis 'Microsoft.Cache/redis@2023-04-01' = {
  location: location
  name: resource_name
  properties: {
    enableNonSslPort: true
    minimumTlsVersion: '1.2'
    sku: {
      capacity: 2
      family: 'C'
      name: 'Standard'
    }
  }
}

