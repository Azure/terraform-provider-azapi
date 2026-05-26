param location string = 'westus2'
param resource_name string = 'acctest0001'

resource replicationFabric2 'Microsoft.RecoveryServices/vaults/replicationFabrics@2022-10-01' = {
  parent: vault
  name: resource_name
  properties: {
    customDetails: {
      instanceType: 'Azure'
      location: location
    }
  }
}

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

