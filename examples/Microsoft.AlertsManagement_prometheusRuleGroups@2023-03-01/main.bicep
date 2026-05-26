param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource account 'Microsoft.Monitor/accounts@2023-04-03' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
  }
}

resource prometheusRuleGroup 'Microsoft.AlertsManagement/prometheusRuleGroups@2023-03-01' = {
  location: location
  name: resource_name
  properties: {
    clusterName: ''
    description: ''
    enabled: false
    rules: [
      {
        enabled: false
        expression: 'histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))\n'
        labels: {
          team: 'prod'
        }
        record: 'job_type:billing_jobs_duration_seconds:99p5m'
      }
    ]
    scopes: [
      account.id
    ]
  }
}

