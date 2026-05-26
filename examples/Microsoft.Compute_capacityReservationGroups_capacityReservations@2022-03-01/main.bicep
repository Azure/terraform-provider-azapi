param location string = 'westus'
param resource_name string = 'acctest0001'

resource capacityReservation 'Microsoft.Compute/capacityReservationGroups/capacityReservations@2022-03-01' = {
  parent: capacityReservationGroup
  location: location
  name: '${resource_name}-ccr'
  sku: {
    capacity: 2
    name: 'Standard_F2'
  }
}

resource capacityReservationGroup 'Microsoft.Compute/capacityReservationGroups@2022-03-01' = {
  location: location
  name: '${resource_name}-ccrg'
}

