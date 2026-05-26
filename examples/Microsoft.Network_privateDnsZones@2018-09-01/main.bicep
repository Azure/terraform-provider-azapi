param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

