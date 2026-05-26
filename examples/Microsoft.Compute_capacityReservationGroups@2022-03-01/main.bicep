param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource capacityReservationGroup 'Microsoft.Compute/capacityReservationGroups@2022-03-01' = {
  location: location
  name: resource_name
}

