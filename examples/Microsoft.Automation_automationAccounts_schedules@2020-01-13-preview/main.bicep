param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource automationAccount 'Microsoft.Automation/automationAccounts@2021-06-22' = {
  location: location
  name: resource_name
  properties: {
    encryption: {
      keySource: 'Microsoft.Automation'
    }
    publicNetworkAccess: true
    sku: {
      name: 'Basic'
    }
  }
}

resource schedule 'Microsoft.Automation/automationAccounts/schedules@2020-01-13-preview' = {
  parent: automationAccount
  name: resource_name
  properties: {
    description: ''
    frequency: 'OneTime'
    startTime: '2024-07-05T08:51:00+00:00'
    timeZone: 'Etc/UTC'
  }
}

