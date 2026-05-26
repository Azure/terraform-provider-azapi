param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.OperationalInsights/clusters@2020-08-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  sku: {
    capacity: 1000
    name: 'CapacityReservation'
  }
}

