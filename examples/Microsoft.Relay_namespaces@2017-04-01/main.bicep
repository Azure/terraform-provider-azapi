param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource namespace 'Microsoft.Relay/namespaces@2017-04-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Standard'
    tier: 'Standard'
  }
}

