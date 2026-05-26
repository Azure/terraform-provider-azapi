param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource spatialAnchorsAccount 'Microsoft.MixedReality/spatialAnchorsAccounts@2021-01-01' = {
  location: location
  name: resource_name
}

