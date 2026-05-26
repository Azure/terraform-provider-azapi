param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource api 'Microsoft.ApiManagement/service/apis@2021-08-01' = {
  parent: service
  name: '${resource_name};rev=1'
  properties: {
    apiType: 'http'
    apiVersion: ''
    format: 'swagger-link-json'
    path: 'test'
    type: 'http'
    value: 'http://conferenceapi.azurewebsites.net/?format=json'
  }
}

resource component 'Microsoft.Insights/components@2020-02-02' = {
  location: location
  name: resource_name
  kind: 'web'
  properties: {
    Application_Type: 'web'
    DisableIpMasking: false
    DisableLocalAuth: false
    ForceCustomerStorageForProfiler: false
    RetentionInDays: 90
    SamplingPercentage: 100
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
  }
}

resource diagnostic 'Microsoft.ApiManagement/service/apis/diagnostics@2021-08-01' = {
  parent: api
  name: 'applicationinsights'
  properties: {
    loggerId: logger.id
    operationNameFormat: 'Name'
  }
}

resource logger 'Microsoft.ApiManagement/service/loggers@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    credentials: {
      instrumentationKey: component.properties.InstrumentationKey
    }
    description: ''
    isBuffered: true
    loggerType: 'applicationInsights'
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

