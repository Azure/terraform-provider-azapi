param location string = 'westeurope'
param resource_name string = 'acctest0001'

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

