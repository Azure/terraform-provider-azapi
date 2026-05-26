param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules@2017-04-01' = {
  parent: notificationHub
  name: resource_name
  properties: {
    rights: [
      'Listen'
    ]
  }
}

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

resource notificationHub 'Microsoft.NotificationHubs/namespaces/notificationHubs@2017-04-01' = {
  parent: namespace
  location: location
  name: resource_name
  properties: {}
}

