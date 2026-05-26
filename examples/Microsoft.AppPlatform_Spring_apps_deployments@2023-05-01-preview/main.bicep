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

resource deployment 'Microsoft.AppPlatform/Spring/apps/deployments@2023-05-01-preview' = {
  parent: app
  name: resource_name
  properties: {
    deploymentSettings: {
      environmentVariables: {}
    }
    source: {
      customContainer: {
        args: []
        command: []
        containerImage: 'springio/gs-spring-boot-docker'
        languageFramework: ''
        server: 'docker.io'
      }
      type: 'Container'
    }
  }
  sku: {
    capacity: 1
    name: 'E0'
    tier: 'Enterprise'
  }
}

