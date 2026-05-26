param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource searchService 'Microsoft.Search/searchServices@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    authOptions: {
      apiKeyOnly: {}
    }
    disableLocalAuth: false
    encryptionWithCmk: {
      enforcement: 'Disabled'
    }
    hostingMode: 'default'
    networkRuleSet: {
      ipRules: []
    }
    partitionCount: 1
    publicNetworkAccess: 'Enabled'
    replicaCount: 1
  }
  sku: {
    name: 'standard'
  }
  tags: {
    environment: 'staging'
  }
}

