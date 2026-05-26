param location string = 'westeurope'
param openid_client_id string = null
param openid_client_secret string = null
param resource_name string = 'acctest0001'

resource openidConnectProvider 'Microsoft.ApiManagement/service/openidConnectProviders@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    clientId: openid_client_id
    clientSecret: openid_client_secret
    description: ''
    displayName: 'Initial Name'
    metadataEndpoint: 'https://azacceptance.hashicorptest.com/example/foo'
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

