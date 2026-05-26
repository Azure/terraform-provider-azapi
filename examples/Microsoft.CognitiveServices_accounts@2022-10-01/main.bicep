param location string = 'westus2'
param resource_name string = 'acctest0001'

resource account 'Microsoft.CognitiveServices/accounts@2022-10-01' = {
  identity: [
    {
      identity_ids: [
        userAssignedIdentity.id
      ]
      type: 'SystemAssigned, UserAssigned'
    }
  ]
  location: location
  name: resource_name
  kind: 'SpeechServices'
  properties: {
    allowedFqdnList: []
    apiProperties: {}
    customSubDomainName: 'acctest-cogacc-230630032807723157'
    disableLocalAuth: false
    dynamicThrottlingEnabled: false
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: false
  }
  sku: {
    name: 'S0'
    tier: 'Standard'
  }
}

resource userAssignedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  location: location
  name: resource_name
}

