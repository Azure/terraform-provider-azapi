param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource resourceProvider 'Microsoft.CustomProviders/resourceProviders@2018-09-01-preview' = {
  location: location
  name: resource_name
  properties: {
    resourceTypes: [
      {
        endpoint: 'https://example.com/'
        name: 'dEf1'
        routingType: 'Proxy'
      }
    ]
  }
}

