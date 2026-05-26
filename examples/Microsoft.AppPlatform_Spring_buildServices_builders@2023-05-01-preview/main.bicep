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

resource builder 'Microsoft.AppPlatform/Spring/buildServices/builders@2023-05-01-preview' = {
  name: resource_name
  properties: {
    buildpackGroups: [
      {
        buildpacks: [
          {
            id: 'tanzu-buildpacks/java-azure'
          }
        ]
        name: 'mix'
      }
    ]
    stack: {
      id: 'io.buildpacks.stacks.bionic'
      version: 'base'
    }
  }
}

