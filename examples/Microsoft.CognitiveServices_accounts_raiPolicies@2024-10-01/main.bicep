param location string = 'eastus'
param resource_name string = 'acctest0003'

resource account 'Microsoft.CognitiveServices/accounts@2022-10-01' = {
  location: location
  name: resource_name
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

resource raiPolicy 'Microsoft.CognitiveServices/accounts/raiPolicies@2024-10-01' = {
  parent: account
  name: 'NoModerationPolicy'
  properties: {
    basePolicyName: 'Microsoft.Default'
    contentFilters: [
      {
        blocking: true
        enabled: true
        name: 'Hate'
        severityThreshold: 'High'
        source: 'Prompt'
      }
    ]
  }
}

