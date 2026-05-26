param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource namespace 'Microsoft.NotificationHubs/namespaces@2017-04-01' = {
  location: location
  name: resource_name
  properties: {
    enabled: true
    namespaceType: 'NotificationHub'
    region: 'westeurope'
  }
  sku: {
    name: 'Free'
  }
}

