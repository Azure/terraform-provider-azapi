param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource afdEndpoint 'Microsoft.Cdn/profiles/afdEndpoints@2021-06-01' = {
  parent: profile
  location: 'global'
  name: resource_name
  properties: {
    enabledState: 'Enabled'
  }
}

resource profile 'Microsoft.Cdn/profiles@2021-06-01' = {
  location: 'global'
  name: resource_name
  properties: {
    originResponseTimeoutSeconds: 120
  }
  sku: {
    name: 'Standard_AzureFrontDoor'
  }
}

