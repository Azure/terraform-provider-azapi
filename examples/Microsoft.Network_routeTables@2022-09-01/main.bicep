param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource routeTable 'Microsoft.Network/routeTables@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    disableBgpRoutePropagation: false
    routes: [
      {
        name: 'first'
        properties: {
          addressPrefix: '10.100.0.0/14'
          nextHopIpAddress: '10.10.1.1'
          nextHopType: 'VirtualAppliance'
        }
      }
    ]
  }
}

