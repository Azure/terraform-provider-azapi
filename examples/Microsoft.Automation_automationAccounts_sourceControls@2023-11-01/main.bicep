param location string = 'westeurope'
param pat string = null
param resource_name string = 'acctest0001'

resource automationAccount 'Microsoft.Automation/automationAccounts@2023-11-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
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

resource sourceControl 'Microsoft.Automation/automationAccounts/sourceControls@2023-11-01' = {
  parent: automationAccount
  name: resource_name
  properties: {
    autoSync: false
    branch: 'master'
    folderPath: '/'
    publishRunbook: false
    repoUrl: 'https://github.com/Azure-Samples/acr-build-helloworld-node.git'
    securityToken: {
      accessToken: pat
      tokenType: 'PersonalAccessToken'
    }
    sourceType: 'GitHub'
  }
}

