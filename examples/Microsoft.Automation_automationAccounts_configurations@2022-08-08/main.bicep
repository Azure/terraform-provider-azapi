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

resource configuration 'Microsoft.Automation/automationAccounts/configurations@2022-08-08' = {
  parent: automationAccount
  location: location
  name: resource_name
  properties: {
    description: 'test'
    logVerbose: false
    source: {
      type: 'embeddedContent'
      value: 'configuration acctest {}'
    }
  }
  tags: {
    ENV: 'prod'
  }
}

