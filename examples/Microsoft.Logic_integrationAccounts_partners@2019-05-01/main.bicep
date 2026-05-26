param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource integrationAccount 'Microsoft.Logic/integrationAccounts@2019-05-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Standard'
  }
}

resource partner 'Microsoft.Logic/integrationAccounts/partners@2019-05-01' = {
  parent: integrationAccount
  name: resource_name
  properties: {
    content: {
      b2b: {
        businessIdentities: [
          {
            qualifier: 'AS2Identity'
            value: 'FabrikamNY'
          }
        ]
      }
    }
    partnerType: 'B2B'
  }
}

