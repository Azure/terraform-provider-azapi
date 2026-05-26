param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource account 'Microsoft.DeviceUpdate/accounts@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    sku: 'Standard'
  }
}

