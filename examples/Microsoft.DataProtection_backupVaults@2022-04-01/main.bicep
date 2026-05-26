param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource backupVault 'Microsoft.DataProtection/backupVaults@2022-04-01' = {
  location: location
  name: resource_name
  properties: {
    storageSettings: [
      {
        datastoreType: 'VaultStore'
        type: 'LocallyRedundant'
      }
    ]
  }
}

