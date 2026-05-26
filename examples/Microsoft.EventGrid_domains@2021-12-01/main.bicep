param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource domain 'Microsoft.EventGrid/domains@2021-12-01' = {
  location: location
  name: resource_name
  properties: {
    autoCreateTopicWithFirstSubscription: true
    autoDeleteTopicWithLastSubscription: true
    disableLocalAuth: false
    inputSchema: 'EventGridSchema'
    inputSchemaMapping: null
    publicNetworkAccess: 'Enabled'
  }
}

