param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource profile 'Microsoft.Cdn/profiles@2021-06-01' = {
  location: 'global'
  name: resource_name
  properties: {
    originResponseTimeoutSeconds: 120
  }
  sku: {
    name: 'Premium_AzureFrontDoor'
  }
}

