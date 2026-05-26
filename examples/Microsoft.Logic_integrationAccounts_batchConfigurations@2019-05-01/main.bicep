param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource batchConfiguration 'Microsoft.Logic/integrationAccounts/batchConfigurations@2019-05-01' = {
  parent: integrationAccount
  name: resource_name
  properties: {
    batchGroupName: 'TestBatchGroup'
    releaseCriteria: {
      messageCount: 80
    }
  }
}

resource integrationAccount 'Microsoft.Logic/integrationAccounts@2019-05-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Standard'
  }
}

