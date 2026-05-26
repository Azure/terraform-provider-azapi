param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ipGroup 'Microsoft.Network/ipGroups@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    ipAddresses: []
  }
  tags: {
    env: 'prod'
  }
}

