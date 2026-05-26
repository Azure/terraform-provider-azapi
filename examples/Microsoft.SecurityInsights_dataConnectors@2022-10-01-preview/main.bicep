param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataConnector 'Microsoft.SecurityInsights/dataConnectors@2022-10-01-preview' = {
  parent: workspace
  name: resource_name
  kind: 'MicrosoftThreatIntelligence'
  properties: {
    dataTypes: {
      bingSafetyPhishingURL: {
        lookbackPeriod: ''
        state: 'Disabled'
      }
      microsoftEmergingThreatFeed: {
        lookbackPeriod: '1970-01-01T00:00:00Z'
        state: 'enabled'
      }
    }
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

resource onboardingState 'Microsoft.SecurityInsights/onboardingStates@2023-06-01-preview' = {
  parent: workspace
  name: 'default'
  properties: {
    customerManagedKey: false
  }
}

resource workspace 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    features: {
      disableLocalAuth: false
      enableLogAccessUsingOnlyResourcePermissions: true
    }
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
    retentionInDays: 30
    sku: {
      name: 'PerGB2018'
    }
    workspaceCapping: {
      dailyQuotaGb: -1
    }
  }
}

