param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource loadTest 'Microsoft.LoadTestService/loadTests@2022-12-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  tags: {
    Team: 'Dev Exp'
  }
  properties: {
    description: 'This is new load test'
  }
}

