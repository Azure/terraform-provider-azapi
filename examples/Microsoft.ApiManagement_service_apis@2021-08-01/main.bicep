param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource api 'Microsoft.ApiManagement/service/apis@2021-08-01' = {
  parent: service
  name: '${resource_name};rev=1'
  properties: {
    apiRevisionDescription: ''
    apiType: 'http'
    apiVersion: ''
    apiVersionDescription: ''
    authenticationSettings: {}
    description: ''
    displayName: 'api1'
    path: 'api1'
    protocols: [
      'https'
    ]
    serviceUrl: ''
    subscriptionRequired: true
    type: 'http'
  }
}

resource service 'Microsoft.ApiManagement/service@2021-08-01' = {
  location: location
  name: resource_name
  properties: {
    certificates: []
    customProperties: {
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11': 'false'
    }
    disableGateway: false
    publicNetworkAccess: 'Enabled'
    publisherEmail: 'pub1@email.com'
    publisherName: 'pub1'
    virtualNetworkType: 'None'
  }
  sku: {
    capacity: 0
    name: 'Consumption'
  }
}

