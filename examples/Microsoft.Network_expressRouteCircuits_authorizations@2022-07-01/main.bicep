param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorization 'Microsoft.Network/expressRouteCircuits/authorizations@2022-07-01' = {
  parent: expressRouteCircuit
  name: resource_name
  properties: {}
}

resource expressRouteCircuit 'Microsoft.Network/expressRouteCircuits@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    authorizationKey: ''
    serviceProviderProperties: {
      bandwidthInMbps: 50
      peeringLocation: 'Silicon Valley'
      serviceProviderName: 'Equinix'
    }
  }
  sku: {
    family: 'MeteredData'
    name: 'Standard_MeteredData'
    tier: 'Standard'
  }
  tags: {
    Environment: 'production'
    Purpose: 'AcceptanceTests'
  }
}

