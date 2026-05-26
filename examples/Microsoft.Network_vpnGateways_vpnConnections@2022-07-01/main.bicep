param location string = 'westeurope'
param resource_name string = 'acctest0001'

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

resource vpnConnection 'Microsoft.Network/vpnGateways/vpnConnections@2022-07-01' = {
  parent: vpnGateway
  name: resource_name
  properties: {
    enableInternetSecurity: false
    remoteVpnSite: {
      id: vpnSite.id
    }
    vpnLinkConnections: [
      {
        name: 'link1'
        properties: {
          connectionBandwidth: 10
          enableBgp: false
          enableRateLimiting: false
          routingWeight: 0
          useLocalAzureIpAddress: false
          usePolicyBasedTrafficSelectors: false
          vpnConnectionProtocolType: 'IKEv2'
          vpnGatewayCustomBgpAddresses: []
          vpnLinkConnectionMode: 'Default'
          vpnSiteLink: {
            id: data.azapi_resource_id.link1.id
          }
        }
      }
      {
        name: 'link2'
        properties: {
          connectionBandwidth: 10
          enableBgp: false
          enableRateLimiting: false
          routingWeight: 0
          useLocalAzureIpAddress: false
          usePolicyBasedTrafficSelectors: false
          vpnConnectionProtocolType: 'IKEv2'
          vpnGatewayCustomBgpAddresses: []
          vpnLinkConnectionMode: 'Default'
          vpnSiteLink: {
            id: data.azapi_resource_id.link2.id
          }
        }
      }
    ]
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

resource vpnSite 'Microsoft.Network/vpnSites@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.1.0/24'
      ]
    }
    virtualWan: {
      id: virtualWan.id
    }
    vpnSiteLinks: [
      {
        name: 'link1'
        properties: {
          fqdn: ''
          ipAddress: '10.0.1.1'
          linkProperties: {
            linkProviderName: ''
            linkSpeedInMbps: 0
          }
        }
      }
      {
        name: 'link2'
        properties: {
          fqdn: ''
          ipAddress: '10.0.1.2'
          linkProperties: {
            linkProviderName: ''
            linkSpeedInMbps: 0
          }
        }
      }
    ]
  }
}

