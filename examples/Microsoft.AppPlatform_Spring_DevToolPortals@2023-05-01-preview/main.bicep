param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource DevToolPortal 'Microsoft.AppPlatform/Spring/DevToolPortals@2023-05-01-preview' = {
  parent: Spring
  name: 'default'
  properties: {
    features: {
      applicationAccelerator: {
        state: 'Disabled'
      }
      applicationLiveView: {
        state: 'Disabled'
      }
    }
    public: false
  }
}

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

