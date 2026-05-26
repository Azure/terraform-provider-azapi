param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ExpressRoutePort 'Microsoft.Network/ExpressRoutePorts@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    bandwidthInGbps: 10
    encapsulation: 'Dot1Q'
    peeringLocation: 'CDC-Canberra'
  }
}

resource expressRouteCircuit 'Microsoft.Network/expressRouteCircuits@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    authorizationKey: ''
    bandwidthInGbps: 5
    expressRoutePort: {
      id: ExpressRoutePort.id
    }
  }
  sku: {
    family: 'MeteredData'
    name: 'Premium_MeteredData'
    tier: 'Premium'
  }
}

