param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dashboard 'Microsoft.Portal/dashboards@2019-01-01-preview' = {
  location: location
  name: resource_name
  properties: {
    lenses: {}
  }
}

