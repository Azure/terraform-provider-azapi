param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.Relay/namespaces/authorizationRules@2017-04-01' = {
  parent: namespace
  name: resource_name
  properties: {
    rights: [
      'Listen'
      'Send'
    ]
  }
}

resource namespace 'Microsoft.Relay/namespaces@2017-04-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Standard'
    tier: 'Standard'
  }
}

