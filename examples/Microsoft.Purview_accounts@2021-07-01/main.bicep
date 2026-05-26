param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource account 'Microsoft.Purview/accounts@2021-07-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
  }
}

