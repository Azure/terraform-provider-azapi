param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.StreamAnalytics/clusters@2020-03-01' = {
  location: location
  name: resource_name
  sku: {
    capacity: 36
    name: 'Default'
  }
}

