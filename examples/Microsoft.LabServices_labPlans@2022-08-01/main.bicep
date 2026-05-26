param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource labPlan 'Microsoft.LabServices/labPlans@2022-08-01' = {
  location: location
  name: resource_name
  properties: {
    allowedRegions: [
      location
    ]
  }
}

