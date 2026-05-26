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

resource powerShell72Module 'Microsoft.Automation/automationAccounts/powerShell72Modules@2020-01-13-preview' = {
  parent: automationAccount
  name: 'xActiveDirectory'
  properties: {
    contentLink: {
      uri: 'https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg'
    }
  }
}

