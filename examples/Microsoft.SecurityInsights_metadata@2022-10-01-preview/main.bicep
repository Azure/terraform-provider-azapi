param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource alertRule 'Microsoft.SecurityInsights/alertRules@2022-10-01-preview' = {
  parent: workspace
  name: resource_name
  kind: 'NRT'
  properties: {
    description: ''
    displayName: 'Some Rule'
    enabled: true
    query: 'AzureActivity |\n  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |\n  where ActivityStatus == "Succeeded" |\n  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller\n'
    severity: 'High'
    suppressionDuration: 'PT5H'
    suppressionEnabled: false
    tactics: []
    techniques: []
  }
}

resource metadata 'Microsoft.SecurityInsights/metadata@2022-10-01-preview' = {
  parent: workspace
  name: resource_name
  properties: {
    contentId: resource_name
    contentSchemaVersion: '2.0'
    kind: 'AnalyticsRule'
    parentId: alertRule.id
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

