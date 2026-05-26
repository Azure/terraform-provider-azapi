param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource backupResourceGuardProxy 'Microsoft.RecoveryServices/vaults/backupResourceGuardProxies@2023-02-01' = {
  parent: vault
  name: resource_name
  properties: {
    resourceGuardResourceId: resourceGuard.id
  }
  type: 'Microsoft.RecoveryServices/vaults/backupResourceGuardProxies'
}

resource resourceGuard 'Microsoft.DataProtection/resourceGuards@2022-04-01' = {
  location: location
  name: resource_name
  properties: {
    vaultCriticalOperationExclusionList: []
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

