param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource component 'Microsoft.Insights/components@2020-02-02' = {
  location: location
  name: resource_name
  kind: 'web'
  properties: {
    Application_Type: 'web'
    DisableIpMasking: false
    DisableLocalAuth: false
    ForceCustomerStorageForProfiler: false
    RetentionInDays: 90
    SamplingPercentage: 100
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
  }
}

resource scheduledQueryRule 'Microsoft.Insights/scheduledQueryRules@2021-08-01' = {
  location: location
  name: resource_name
  kind: 'LogAlert'
  properties: {
    autoMitigate: false
    checkWorkspaceAlertsStorageConfigured: false
    criteria: {
      allOf: [
        {
          dimensions: null
          operator: 'Equal'
          query: ' requests\n| summarize CountByCountry=count() by client_CountryOrRegion\n'
          threshold: 5
          timeAggregation: 'Count'
        }
      ]
    }
    enabled: true
    evaluationFrequency: 'PT5M'
    scopes: [
      component.id
    ]
    severity: 3
    skipQueryValidation: false
    targetResourceTypes: null
    windowSize: 'PT5M'
  }
}

