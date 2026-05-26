param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource availabilitySet 'Microsoft.Compute/availabilitySets@2021-11-01' = {
  location: location
  name: resource_name
  properties: {
    platformFaultDomainCount: 3
    platformUpdateDomainCount: 5
  }
  sku: {
    name: 'Aligned'
  }
}

