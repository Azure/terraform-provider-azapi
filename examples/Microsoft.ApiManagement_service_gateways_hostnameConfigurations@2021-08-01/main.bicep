param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource certificate 'Microsoft.ApiManagement/service/certificates@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    data: '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'
    password: ''
  }
}

resource gateway 'Microsoft.ApiManagement/service/gateways@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    description: ''
    locationData: {
      city: ''
      countryOrRegion: ''
      district: ''
      name: 'test'
    }
  }
}

resource hostnameConfiguration 'Microsoft.ApiManagement/service/gateways/hostnameConfigurations@2021-08-01' = {
  parent: gateway
  name: resource_name
  properties: {
    certificateId: certificate.id
    hostname: 'host-name-230630032559671742'
    http2Enabled: true
    negotiateClientCertificate: false
    tls10Enabled: false
    tls11Enabled: false
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
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA256': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_GCM_SHA256': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA256': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_GCM_SHA384': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30': 'false'
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
    capacity: 1
    name: 'Developer'
  }
}

