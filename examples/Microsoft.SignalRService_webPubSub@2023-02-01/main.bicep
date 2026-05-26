param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource webPubSub 'Microsoft.SignalRService/webPubSub@2023-02-01' = {
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

