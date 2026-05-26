param location string = 'eastus'
param resource_name string = 'acctest0001'

resource dataBoxEdgeDevice 'Microsoft.DataBoxEdge/dataBoxEdgeDevices@2022-03-01' = {
  location: location
  name: resource_name
  sku: {
    name: 'EdgeP_Base'
    tier: 'Standard'
  }
}

resource mobileNetwork 'Microsoft.MobileNetwork/mobileNetworks@2022-11-01' = {
  location: location
  name: resource_name
  properties: {
    publicLandMobileNetworkIdentifier: {
      mcc: '001'
      mnc: '01'
    }
  }
}

resource packetCoreControlPlane 'Microsoft.MobileNetwork/packetCoreControlPlanes@2022-11-01' = {
  location: location
  name: resource_name
  properties: {
    controlPlaneAccessInterface: {}
    localDiagnosticsAccess: {
      authenticationType: 'AAD'
    }
    platform: {
      azureStackEdgeDevice: {
        id: dataBoxEdgeDevice.id
      }
      type: 'AKS-HCI'
    }
    sites: [
      {
        id: site.id
      }
    ]
    sku: 'G0'
    ueMtu: 1440
  }
}

resource packetCoreDataPlane 'Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreDataPlanes@2022-11-01' = {
  parent: packetCoreControlPlane
  location: location
  name: resource_name
  properties: {
    userPlaneAccessInterface: {}
  }
}

resource site 'Microsoft.MobileNetwork/mobileNetworks/sites@2022-11-01' = {
  parent: mobileNetwork
  location: location
  name: resource_name
  properties: {}
}

