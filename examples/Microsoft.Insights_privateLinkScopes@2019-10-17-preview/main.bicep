param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource privateLinkScope 'Microsoft.Insights/privateLinkScopes@2019-10-17-preview' = {
  location: 'Global'
  name: resource_name
  properties: {}
}

