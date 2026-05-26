param express_route_shared_key string = null
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

resource peering 'Microsoft.Network/expressRouteCircuits/peerings@2022-07-01' = {
  parent: expressRouteCircuit
  name: 'AzurePrivatePeering'
  properties: {
    azureASN: 12076
    gatewayManagerEtag: ''
    peerASN: 100
    peeringType: 'AzurePrivatePeering'
    primaryPeerAddressPrefix: '192.168.1.0/30'
    secondaryPeerAddressPrefix: '192.168.2.0/30'
    sharedKey: express_route_shared_key
    state: 'Enabled'
    vlanId: 100
  }
}

