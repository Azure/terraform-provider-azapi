param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource fhirService 'Microsoft.HealthcareApis/workspaces/fhirServices@2022-12-01' = {
  parent: workspace
  location: location
  name: resource_name
  kind: 'fhir-R4'
  properties: {
    acrConfiguration: {}
    authenticationConfiguration: {
      audience: 'https://acctestfhir.fhir.azurehealthcareapis.com'
      authority: 'https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}'
      smartProxyEnabled: false
    }
    corsConfiguration: {
      allowCredentials: false
      headers: []
      methods: []
      origins: []
    }
  }
}

resource fhirService2 'Microsoft.HealthcareApis/workspaces/fhirServices@2022-12-01' = {
  parent: workspace
  location: location
  name: resource_name
  kind: 'fhir-R4'
  properties: {
    acrConfiguration: {}
    authenticationConfiguration: {
      audience: fhirService.properties.authenticationConfiguration.audience
      authority: fhirService.properties.authenticationConfiguration.authority
      smartProxyEnabled: false
    }
    corsConfiguration: {
      allowCredentials: false
      headers: []
      methods: []
      origins: []
    }
  }
}

resource workspace 'Microsoft.HealthcareApis/workspaces@2022-12-01' = {
  location: location
  name: resource_name
}

