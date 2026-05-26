param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource proximityPlacementGroup 'Microsoft.Compute/proximityPlacementGroups@2022-03-01' = {
  location: location
  name: resource_name
  properties: {}
}

