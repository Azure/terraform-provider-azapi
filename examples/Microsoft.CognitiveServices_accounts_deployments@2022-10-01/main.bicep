param location string = 'eastus'
param resource_name string = 'acctest0003'

resource account 'Microsoft.CognitiveServices/accounts@2022-10-01' = {
  location: location
  name: resource_name
  identity: {
    type: 'None'
    userAssignedIdentities: null
  }
  kind: 'OpenAI'
  properties: {
    disableLocalAuth: false
    dynamicThrottlingEnabled: false
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: false
  }
  sku: {
    name: 'S0'
  }
}

resource deployment 'Microsoft.CognitiveServices/accounts/deployments@2023-05-01' = {
  parent: account
  name: 'testdep'
  properties: {
    model: {
      format: 'OpenAI'
      name: 'text-embedding-ada-002'
    }
  }
}

