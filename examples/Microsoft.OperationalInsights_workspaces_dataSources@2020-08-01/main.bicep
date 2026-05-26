param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataSource 'Microsoft.OperationalInsights/workspaces/dataSources@2020-08-01' = {
  parent: workspace
  name: resource_name
  kind: 'WindowsPerformanceCounter'
  properties: {
    counterName: 'CPU'
    instanceName: '*'
    intervalSeconds: 10
    objectName: 'CPU'
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

