param location string = 'eastus'
param resource_name string = 'acctest0001'

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

