param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource AAAA 'Microsoft.Network/dnsZones/AAAA@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    AAAARecords: [
      {
        ipv6Address: '2607:f8b0:4009:1803::1005'
      }
      {
        ipv6Address: '2607:f8b0:4009:1803::1006'
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

