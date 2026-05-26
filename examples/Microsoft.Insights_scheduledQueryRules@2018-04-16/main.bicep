param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource scheduledQueryRule 'Microsoft.Insights/scheduledQueryRules@2018-04-16' = {
  location: location
  name: resource_name
  properties: {
    action: {
      criteria: [
        {
          dimensions: [
            {
              name: 'InstanceName'
              operator: 'Include'
              values: [
                '1'
              ]
            }
          ]
          metricName: 'Average_% Idle Time'
        }
      ]
      'odata.type': 'Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.LogToMetricAction'
    }
    description: ''
    enabled: 'true'
    source: {
      authorizedResources: []
      dataSourceId: workspace.id
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

