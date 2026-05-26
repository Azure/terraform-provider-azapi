param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource hub 'Microsoft.SignalRService/webPubSub/hubs@2023-02-01' = {
  parent: webPubSub
  name: resource_name
  properties: {
    anonymousConnectPolicy: 'Deny'
    eventListeners: []
  }
}

resource webPubSub 'Microsoft.SignalRService/webPubSub@2023-02-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    disableAadAuth: false
    disableLocalAuth: false
    publicNetworkAccess: 'Enabled'
    tls: {
      clientCertEnabled: false
    }
  }
  sku: {
    capacity: 1
    name: 'Standard_S1'
  }
}

