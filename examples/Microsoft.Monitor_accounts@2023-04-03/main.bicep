param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource account 'Microsoft.Monitor/accounts@2023-04-03' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
  }
}

