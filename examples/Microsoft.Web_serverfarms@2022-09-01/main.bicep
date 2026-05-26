param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource serverfarm 'Microsoft.Web/serverfarms@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    hyperV: false
    perSiteScaling: false
    reserved: false
    zoneRedundant: false
  }
  sku: {
    name: 'S1'
  }
}

