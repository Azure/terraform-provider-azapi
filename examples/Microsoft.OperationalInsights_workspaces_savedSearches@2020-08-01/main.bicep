param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource savedSearch 'Microsoft.OperationalInsights/workspaces/savedSearches@2020-08-01' = {
  parent: workspace
  name: resource_name
  properties: {
    category: 'Saved Search Test Category'
    displayName: 'Create or Update Saved Search Test'
    functionAlias: ''
    query: 'Heartbeat | summarize Count() by Computer | take a'
    tags: []
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

