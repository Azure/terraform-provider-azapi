param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource routeFilter 'Microsoft.Network/routeFilters@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    rules: []
  }
}

