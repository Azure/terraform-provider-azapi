param express_route_connection_shared_key string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ExpressRoutePort 'Microsoft.Network/ExpressRoutePorts@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    bandwidthInGbps: 10
    encapsulation: 'Dot1Q'
    peeringLocation: 'Airtel-Chennai2-CLS'
  }
}

resource ExpressRoutePort2 'Microsoft.Network/ExpressRoutePorts@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    bandwidthInGbps: 10
    encapsulation: 'Dot1Q'
    peeringLocation: 'CDC-Canberra'
  }
}

resource connection 'Microsoft.Network/expressRouteCircuits/peerings/connections@2022-07-01' = {
  parent: peering
  name: resource_name
  properties: {
    addressPrefix: '192.169.8.0/29'
    expressRouteCircuitPeering: {
      id: peering.id
    }
    peerExpressRouteCircuitPeering: {
      id: peering2.id
    }
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
    name: 'Standard_MeteredData'
    tier: 'Standard'
  }
}

resource expressRouteCircuit2 'Microsoft.Network/expressRouteCircuits@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    authorizationKey: ''
    bandwidthInGbps: 5
    expressRoutePort: {
      id: ExpressRoutePort2.id
    }
  }
  sku: {
    family: 'MeteredData'
    name: 'Standard_MeteredData'
    tier: 'Standard'
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
    secondaryPeerAddressPrefix: '192.168.1.0/30'
    sharedKey: express_route_connection_shared_key
    state: 'Enabled'
    vlanId: 100
  }
}

resource peering2 'Microsoft.Network/expressRouteCircuits/peerings@2022-07-01' = {
  parent: expressRouteCircuit2
  name: 'AzurePrivatePeering'
  properties: {
    azureASN: 12076
    gatewayManagerEtag: ''
    peerASN: 100
    peeringType: 'AzurePrivatePeering'
    primaryPeerAddressPrefix: '192.168.1.0/30'
    secondaryPeerAddressPrefix: '192.168.1.0/30'
    sharedKey: express_route_connection_shared_key
    state: 'Enabled'
    vlanId: 100
  }
}

