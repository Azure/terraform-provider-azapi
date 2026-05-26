param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource replicationPolicy 'Microsoft.RecoveryServices/vaults/replicationPolicies@2022-10-01' = {
  parent: vault
  name: resource_name
  properties: {
    providerSpecificInput: {
      appConsistentFrequencyInMinutes: 240
      crashConsistentFrequencyInMinutes: 10
      enableMultiVmSync: 'True'
      instanceType: 'InMageRcm'
      recoveryPointHistoryInMinutes: 1440
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

