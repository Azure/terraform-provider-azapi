param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ddosProtectionPlan 'Microsoft.Network/ddosProtectionPlans@2022-07-01' = {
  location: location
  name: resource_name
}

