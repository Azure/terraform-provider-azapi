param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ExpressRoutePort 'Microsoft.Network/ExpressRoutePorts@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    bandwidthInGbps: 10
    billingType: 'MeteredData'
    encapsulation: 'Dot1Q'
    peeringLocation: 'Airtel-Chennai2-CLS'
  }
  tags: {
    ENV: 'Test'
  }
}

resource authorization 'Microsoft.Network/expressRoutePorts/authorizations@2022-07-01' = {
  parent: ExpressRoutePort
  name: resource_name
  properties: {}
}

