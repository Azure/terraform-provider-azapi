param location string = 'westus2'
param resource_name string = 'acctest0001'

resource service 'Microsoft.HealthcareApis/services@2022-12-01' = {
  location: location
  name: resource_name
  kind: 'fhir'
  properties: {
    accessPolicies: [
      {
        objectId: data.azurerm_client_config.current.object_id
      }
    ]
    authenticationConfiguration: {}
    corsConfiguration: {}
    cosmosDbConfiguration: {
      offerThroughput: 1000
    }
    publicNetworkAccess: 'Enabled'
  }
}

