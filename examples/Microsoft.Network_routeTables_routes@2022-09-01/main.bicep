param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource route 'Microsoft.Network/routeTables/routes@2022-09-01' = {
  parent: routeTable
  name: resource_name
  properties: {
    addressPrefix: '10.1.0.0/16'
    nextHopType: 'VnetLocal'
  }
}

resource routeTable 'Microsoft.Network/routeTables@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    disableBgpRoutePropagation: false
  }
}

