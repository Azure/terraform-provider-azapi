param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource vault 'Microsoft.RecoveryServices/vaults@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    name: 'Standard'
  }
}

