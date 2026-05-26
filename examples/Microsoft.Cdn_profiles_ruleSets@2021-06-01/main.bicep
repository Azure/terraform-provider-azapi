param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource profile 'Microsoft.Cdn/profiles@2021-06-01' = {
  location: 'global'
  name: resource_name
  properties: {
    originResponseTimeoutSeconds: 120
  }
  sku: {
    name: 'Standard_AzureFrontDoor'
  }
}

resource ruleSet 'Microsoft.Cdn/profiles/ruleSets@2021-06-01' = {
  parent: profile
  name: resource_name
}

