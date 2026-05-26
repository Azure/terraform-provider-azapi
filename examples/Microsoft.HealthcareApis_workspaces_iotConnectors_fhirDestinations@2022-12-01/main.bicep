param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource consumerGroup 'Microsoft.EventHub/namespaces/eventhubs/consumerGroups@2021-11-01' = {
  parent: eventhub
  name: resource_name
  properties: {
    userMetadata: ''
  }
}

resource eventhub 'Microsoft.EventHub/namespaces/eventhubs@2021-11-01' = {
  parent: namespace
  name: resource_name
  properties: {
    messageRetentionInDays: 1
    partitionCount: 2
    status: 'Active'
  }
}

resource fhirDestination 'Microsoft.HealthcareApis/workspaces/iotConnectors/fhirDestinations@2022-12-01' = {
  parent: iotConnector
  location: location
  name: resource_name
  properties: {
    fhirMapping: {
      content: {
        template: []
        templateType: 'CollectionFhirTemplate'
      }
    }
    fhirServiceResourceId: fhirService.id
    resourceIdentityResolutionType: 'Create'
  }
}

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

resource iotConnector 'Microsoft.HealthcareApis/workspaces/iotConnectors@2022-12-01' = {
  parent: workspace
  location: location
  name: resource_name
  properties: {
    deviceMapping: {
      content: {
        template: []
        templateType: 'CollectionContent'
      }
    }
    ingestionEndpointConfiguration: {
      consumerGroup: consumerGroup.id
      eventHubName: eventhub.name
      fullyQualifiedEventHubNamespace: '${namespace.name}.servicebus.windows.net'
    }
  }
}

resource namespace 'Microsoft.EventHub/namespaces@2022-01-01-preview' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    isAutoInflateEnabled: false
    publicNetworkAccess: 'Enabled'
    zoneRedundant: false
  }
  sku: {
    capacity: 1
    name: 'Standard'
    tier: 'Standard'
  }
}

resource workspace 'Microsoft.HealthcareApis/workspaces@2022-12-01' = {
  location: location
  name: resource_name
}

