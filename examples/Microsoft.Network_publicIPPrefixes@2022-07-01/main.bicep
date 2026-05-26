param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource publicIPPrefix 'Microsoft.Network/publicIPPrefixes@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    prefixLength: 30
    publicIPAddressVersion: 'IPv4'
  }
  sku: {
    name: 'Standard'
  }
  zones: [
    '1'
  ]
}

