param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource A 'Microsoft.Network/privateDnsZones/A@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    aRecords: [
      {
        ipv4Address: '1.2.4.5'
      }
      {
        ipv4Address: '1.2.3.4'
      }
    ]
    metadata: {}
    ttl: 300
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

