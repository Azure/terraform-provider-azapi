param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource queryPack 'Microsoft.OperationalInsights/queryPacks@2019-09-01' = {
  location: location
  name: resource_name
  properties: {}
}

