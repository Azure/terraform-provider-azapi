param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.EventHub/clusters@2021-11-01' = {
  location: location
  name: resource_name
  sku: {
    capacity: 1
    name: 'Dedicated'
  }
}

