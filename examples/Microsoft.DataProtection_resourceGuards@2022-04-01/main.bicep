param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource resourceGuard 'Microsoft.DataProtection/resourceGuards@2022-04-01' = {
  location: location
  name: resource_name
  properties: {
    vaultCriticalOperationExclusionList: []
  }
}

