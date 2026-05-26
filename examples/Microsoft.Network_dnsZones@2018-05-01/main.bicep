param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

