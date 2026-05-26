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

resource applicationAccelerator 'Microsoft.AppPlatform/Spring/applicationAccelerators@2023-05-01-preview' = {
  parent: Spring
  name: 'default'
}

resource customizedAccelerator 'Microsoft.AppPlatform/Spring/applicationAccelerators/customizedAccelerators@2023-05-01-preview' = {
  parent: applicationAccelerator
  name: resource_name
  properties: {
    description: ''
    displayName: ''
    gitRepository: {
      authSetting: {
        authType: 'Public'
      }
      branch: 'master'
      commit: ''
      gitTag: ''
      url: 'https://github.com/Azure-Samples/piggymetrics'
    }
    iconUrl: ''
  }
}

