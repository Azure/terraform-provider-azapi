param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource natRule 'Microsoft.Network/vpnGateways/natRules@2022-07-01' = {
  parent: vpnGateway
  name: resource_name
  properties: {
    externalMappings: [
      {
        addressSpace: '192.168.21.0/26'
      }
    ]
    internalMappings: [
      {
        addressSpace: '10.4.0.0/26'
      }
    ]
    mode: 'EgressSnat'
    type: 'Static'
  }
}

resource virtualHub 'Microsoft.Network/virtualHubs@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressPrefix: '10.0.0.0/24'
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

resource vpnGateway 'Microsoft.Network/vpnGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    enableBgpRouteTranslationForNat: false
    isRoutingPreferenceInternet: false
    virtualHub: {
      id: virtualHub.id
    }
    vpnGatewayScaleUnit: 1
  }
}

