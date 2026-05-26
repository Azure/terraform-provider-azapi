param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource signalR 'Microsoft.SignalRService/signalR@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    cors: {}
    disableAadAuth: false
    disableLocalAuth: false
    features: [
      {
        flag: 'ServiceMode'
        value: 'Default'
      }
      {
        flag: 'EnableConnectivityLogs'
        value: 'False'
      }
      {
        flag: 'EnableMessagingLogs'
        value: 'False'
      }
      {
        flag: 'EnableLiveTrace'
        value: 'False'
      }
    ]
    publicNetworkAccess: 'Enabled'
    resourceLogConfiguration: {
      categories: [
        {
          enabled: 'false'
          name: 'MessagingLogs'
        }
        {
          enabled: 'false'
          name: 'ConnectivityLogs'
        }
        {
          enabled: 'false'
          name: 'HttpRequestLogs'
        }
      ]
    }
    serverless: {
      connectionTimeoutInSeconds: 30
    }
    tls: {
      clientCertEnabled: false
    }
    upstream: {
      templates: []
    }
  }
  sku: {
    capacity: 1
    name: 'Standard_S1'
  }
}

