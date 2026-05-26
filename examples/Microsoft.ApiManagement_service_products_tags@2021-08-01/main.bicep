param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource product 'Microsoft.ApiManagement/service/products@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    description: ''
    displayName: 'Test Product'
    state: 'notPublished'
    subscriptionRequired: false
    terms: ''
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

resource service_tag 'Microsoft.ApiManagement/service/tags@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    displayName: resource_name
  }
}

resource tag 'Microsoft.ApiManagement/service/products/tags@2021-08-01' = {
  parent: product
  name: service_tag.name
}

