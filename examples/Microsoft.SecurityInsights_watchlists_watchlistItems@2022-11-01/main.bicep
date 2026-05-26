param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource onboardingState 'Microsoft.SecurityInsights/onboardingStates@2022-11-01' = {
  parent: workspace
  name: 'default'
  properties: {
    customerManagedKey: false
  }
}

resource watchlist 'Microsoft.SecurityInsights/watchlists@2022-11-01' = {
  parent: workspace
  name: resource_name
  properties: {
    displayName: 'test'
    itemsSearchKey: 'k1'
    provider: 'Microsoft'
    source: ''
  }
}

resource watchlistItem 'Microsoft.SecurityInsights/watchlists/watchlistItems@2022-11-01' = {
  parent: watchlist
  name: '196abd06-eb4e-4322-9c70-37c32e1a588a'
  properties: {
    itemsKeyValue: {
      k1: 'v1'
    }
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

