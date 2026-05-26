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

resource domain 'Microsoft.AppPlatform/Spring/gateways/domains@2023-05-01-preview' = {
  parent: gateway
  name: '${resource_name}.azuremicroservices.io'
  properties: {
    thumbprint: ''
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

