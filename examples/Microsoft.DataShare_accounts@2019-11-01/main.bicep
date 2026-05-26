param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource account 'Microsoft.DataShare/accounts@2019-11-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  tags: {
    env: 'Test'
  }
}

