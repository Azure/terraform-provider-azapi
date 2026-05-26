param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource hostGroup 'Microsoft.Compute/hostGroups@2021-11-01' = {
  location: location
  name: resource_name
  properties: {
    platformFaultDomainCount: 2
  }
}

