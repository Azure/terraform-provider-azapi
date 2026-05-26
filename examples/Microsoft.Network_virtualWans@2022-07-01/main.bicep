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

