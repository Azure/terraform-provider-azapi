param location string = 'westeurope'
param resource_name string = 'acctest0001'
param service_principal_application_id string = null

resource automationAccount 'Microsoft.Automation/automationAccounts@2021-06-22' = {
  location: location
  name: resource_name
  properties: {
    encryption: {
      keySource: 'Microsoft.Automation'
    }
    publicNetworkAccess: true
    sku: {
      name: 'Basic'
    }
  }
}

resource connection 'Microsoft.Automation/automationAccounts/connections@2020-01-13-preview' = {
  parent: automationAccount
  name: resource_name
  properties: {
    connectionType: {
      name: 'AzureServicePrincipal'
    }
    description: ''
    fieldDefinitionValues: {
      ApplicationId: service_principal_application_id
      CertificateThumbprint: 'AEB97B81A68E8988850972916A8B8B6CD8F39813\n'
      SubscriptionId: data.azurerm_client_config.current.subscription_id
      TenantId: data.azurerm_client_config.current.tenant_id
    }
  }
}

