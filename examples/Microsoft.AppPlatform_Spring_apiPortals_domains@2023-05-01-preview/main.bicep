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

resource apiPortal 'Microsoft.AppPlatform/Spring/apiPortals@2023-05-01-preview' = {
  parent: Spring
  name: 'default'
  properties: {
    gatewayIds: []
    httpsOnly: false
    public: false
  }
  sku: {
    capacity: 1
    name: 'E0'
    tier: 'Enterprise'
  }
}

resource domain 'Microsoft.AppPlatform/Spring/apiPortals/domains@2023-05-01-preview' = {
  parent: apiPortal
  name: '${resource_name}.azuremicroservices.io'
  properties: {
    thumbprint: ''
  }
}

