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
    displayName: 'api1'
    path: 'api1'
    protocols: [
      'https'
    ]
    subscriptionRequired: true
    type: 'http'
  }
}

resource service 'Microsoft.ApiManagement/service@2022-08-01' = {
  location: location
  name: '${resource_name}-service'
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

resource tag 'Microsoft.ApiManagement/service/tags@2022-08-01' = {
  parent: service
  name: '${resource_name}-tag'
  properties: {
    displayName: '${resource_name}-tag'
  }
}

resource tagDescription 'Microsoft.ApiManagement/service/apis/tagDescriptions@2022-08-01' = {
  parent: api
  name: '${resource_name}-tag'
  properties: {
    description: 'tag description'
    externalDocsDescription: 'external tag description'
    externalDocsUrl: 'https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs'
  }
}

resource tag_1 'Microsoft.ApiManagement/service/apis/tags@2022-08-01' = {
  parent: api
  name: '${resource_name}-tag'
}

