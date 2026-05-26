param location string = 'westus'
param resource_name string = 'acctest0001'

resource account 'Microsoft.CognitiveServices/accounts@2024-10-01' = {
  location: location
  name: '${resource_name}-ca'
  kind: 'OpenAI'
  properties: {
    allowedFqdnList: []
    apiProperties: {}
    customSubDomainName: ''
    disableLocalAuth: false
    dynamicThrottlingEnabled: false
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: false
  }
  sku: {
    name: 'S0'
  }
}

resource raiBlocklist 'Microsoft.CognitiveServices/accounts/raiBlocklists@2024-10-01' = {
  parent: account
  name: '${resource_name}-crb'
  properties: {
    description: 'Acceptance test data new azurerm resource'
  }
}

