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

resource applicationLiveView 'Microsoft.AppPlatform/Spring/applicationLiveViews@2023-05-01-preview' = {
  parent: Spring
  name: 'default'
}

