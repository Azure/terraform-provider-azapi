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

resource share 'Microsoft.DataShare/accounts/shares@2019-11-01' = {
  parent: account
  name: resource_name
  properties: {
    description: ''
    shareKind: 'CopyBased'
    terms: ''
  }
}

