param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource topic 'Microsoft.EventGrid/topics@2021-12-01' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    inputSchema: 'EventGridSchema'
    inputSchemaMapping: null
    publicNetworkAccess: 'Enabled'
  }
}

