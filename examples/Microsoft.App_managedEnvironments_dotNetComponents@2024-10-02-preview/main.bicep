param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource aspireDashboard 'Microsoft.App/managedEnvironments/dotNetComponents@2024-10-02-preview' = {
  parent: managedEnvironment
  name: resource_name
  properties: {
    componentType: 'AspireDashboard'
    configurations: []
    serviceBinds: []
  }
}

resource managedEnvironment 'Microsoft.App/managedEnvironments@2022-03-01' = {
  location: location
  name: resource_name
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: workspace.properties.customerId
        sharedKey: data.azapi_resource_action.sharedKeys.output.primarySharedKey
      }
    }
    vnetConfiguration: {}
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

