param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource natGateway 'Microsoft.Network/natGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    idleTimeoutInMinutes: 10
  }
  sku: {
    name: 'Standard'
  }
}

