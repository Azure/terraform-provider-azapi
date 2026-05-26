param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource staticSite 'Microsoft.Web/staticSites@2021-02-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Free'
    tier: 'Free'
  }
}

