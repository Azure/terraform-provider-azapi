param location string = 'westus'
param resource_name string = 'acctest0001'

resource api 'Microsoft.ApiManagement/service/apis@2022-08-01' = {
  parent: service
  name: '${resource_name}-api;rev=1'
  properties: {
    apiRevisionDescription: ''
    apiType: 'http'
    apiVersionDescription: ''
    authenticationSettings: {}
    description: 'What is my purpose? You parse butter.'
    displayName: 'Butter Parser'
    path: 'butter-parser'
    protocols: [
      'http'
      'https'
    ]
    serviceUrl: 'https://example.com/foo/bar'
    subscriptionKeyParameterNames: {
      header: 'X-Butter-Robot-API-Key'
      query: 'location'
    }
    subscriptionRequired: true
    type: 'http'
  }
}

resource operation 'Microsoft.ApiManagement/service/apis/operations@2022-08-01' = {
  parent: api
  name: '${resource_name}-operation'
  properties: {
    description: ''
    displayName: 'DELETE Resource'
    method: 'DELETE'
    responses: []
    templateParameters: []
    urlTemplate: '/resource'
  }
}

resource service 'Microsoft.ApiManagement/service@2022-08-01' = {
  location: location
  name: '${resource_name}-am'
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

