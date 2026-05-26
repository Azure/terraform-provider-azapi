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

resource session 'Microsoft.Logic/integrationAccounts/sessions@2019-05-01' = {
  parent: integrationAccount
  name: resource_name
  properties: {
    content: '\t{\n       "controlNumber": "1234"\n    }\n'
  }
}

