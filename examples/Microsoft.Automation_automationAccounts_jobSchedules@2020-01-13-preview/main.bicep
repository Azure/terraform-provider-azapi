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

resource jobSchedule 'Microsoft.Automation/automationAccounts/jobSchedules@2020-01-13-preview' = {
  parent: automationAccount
  name: '194a324f-9e3d-43ee-1234-c968b797edd5'
  properties: {
    runbook: {
      name: runbook.name
    }
    schedule: {
      name: schedule.name
    }
  }
}

resource runbook 'Microsoft.Automation/automationAccounts/runbooks@2019-06-01' = {
  parent: automationAccount
  location: location
  name: 'Get-AzureVMTutorial'
  properties: {
    description: 'This is a test runbook for terraform acceptance test'
    draft: {}
    logActivityTrace: 0
    logProgress: true
    logVerbose: true
    runbookType: 'PowerShell'
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

