param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource batchAccount 'Microsoft.Batch/batchAccounts@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    encryption: {
      keySource: 'Microsoft.Batch'
    }
    poolAllocationMode: 'BatchService'
    publicNetworkAccess: 'Enabled'
  }
}

