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

