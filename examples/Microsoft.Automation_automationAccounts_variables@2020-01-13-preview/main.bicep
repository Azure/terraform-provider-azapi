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

resource variable 'Microsoft.Automation/automationAccounts/variables@2020-01-13-preview' = {
  parent: automationAccount
  name: resource_name
  properties: {
    description: ''
    isEncrypted: false
    value: '"Hello, Terraform Basic Test."'
  }
}

