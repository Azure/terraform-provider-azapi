param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource integrationAccount 'Microsoft.Logic/integrationAccounts@2019-05-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Basic'
  }
}

