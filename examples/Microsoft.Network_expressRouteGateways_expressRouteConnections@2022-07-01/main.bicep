param location string = 'westeurope'
param resource_name string = 'acctest0001'
param shared_key string = null

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

resource expressRouteConnection 'Microsoft.Network/expressRouteGateways/expressRouteConnections@2022-07-01' = {
  parent: expressRouteGateway
  name: resource_name
  properties: {
    enableInternetSecurity: false
    expressRouteCircuitPeering: {
      id: peering.id
    }
    expressRouteGatewayBypass: false
    routingConfiguration: {}
    routingWeight: 0
  }
}

resource expressRouteGateway 'Microsoft.Network/expressRouteGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    allowNonVirtualWanTraffic: false
    autoScaleConfiguration: {
      bounds: {
        min: 1
      }
    }
    virtualHub: {
      id: virtualHub.id
    }
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
    sharedKey: shared_key
    state: 'Enabled'
    vlanId: 100
  }
}

resource virtualHub 'Microsoft.Network/virtualHubs@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressPrefix: '10.0.1.0/24'
    hubRoutingPreference: 'ExpressRoute'
    virtualRouterAutoScaleConfiguration: {
      minCapacity: 2
    }
    virtualWan: {
      id: virtualWan.id
    }
  }
}

resource virtualWan 'Microsoft.Network/virtualWans@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    allowBranchToBranchTraffic: true
    disableVpnEncryption: false
    office365LocalBreakoutCategory: 'None'
    type: 'Standard'
  }
}

