param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource Spring 'Microsoft.AppPlatform/Spring@2023-05-01-preview' = {
  location: location
  name: resource_name
  properties: {
    zoneRedundant: false
  }
  sku: {
    name: 'E0'
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

resource gateway 'Microsoft.AppPlatform/Spring/gateways@2023-05-01-preview' = {
  parent: Spring
  name: 'default'
  properties: {
    httpsOnly: false
    public: false
  }
  sku: {
    capacity: 1
    name: 'E0'
    tier: 'Enterprise'
  }
}

resource routeConfig 'Microsoft.AppPlatform/Spring/gateways/routeConfigs@2023-05-01-preview' = {
  parent: gateway
  name: resource_name
  properties: {
    appResourceId: app.id
    protocol: 'HTTP'
    routes: []
    ssoEnabled: false
  }
}

