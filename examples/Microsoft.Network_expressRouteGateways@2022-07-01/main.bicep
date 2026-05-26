param location string = 'westeurope'
param resource_name string = 'acctest0001'

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

