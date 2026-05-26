param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cache 'Microsoft.ApiManagement/service/caches@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    connectionString: '${redis.name}.redis.cache.windows.net:6380,password=${data.azapi_resource_action.listKeys.output.primaryKey},ssl=true,abortConnect=False'
    useFromLocation: 'default'
  }
}

resource redis 'Microsoft.Cache/redis@2023-04-01' = {
  location: 'eastus'
  name: resource_name
  properties: {
    enableNonSslPort: true
    minimumTlsVersion: '1.2'
    sku: {
      capacity: 2
      family: 'C'
      name: 'Standard'
    }
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

