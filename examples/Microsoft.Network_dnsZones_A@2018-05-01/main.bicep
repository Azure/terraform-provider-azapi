param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource A 'Microsoft.Network/dnsZones/A@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    ARecords: [
      {
        ipv4Address: '1.2.4.5'
      }
      {
        ipv4Address: '1.2.3.4'
      }
    ]
    TTL: 300
    metadata: {}
    targetResource: {}
  }
}

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

