param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource automationRule 'Microsoft.SecurityInsights/automationRules@2022-10-01-preview' = {
  parent: workspace
  name: '3b862818-ad7b-409e-83be-8812f2a06d37'
  properties: {
    actions: [
      {
        actionConfiguration: {
          classification: ''
          classificationComment: ''
          classificationReason: ''
          severity: ''
          status: 'Active'
        }
        actionType: 'ModifyProperties'
        order: 1
      }
    ]
    displayName: 'acctest-SentinelAutoRule-230630033910945846'
    order: 1
    triggeringLogic: {
      isEnabled: true
      triggersOn: 'Incidents'
      triggersWhen: 'Created'
    }
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

