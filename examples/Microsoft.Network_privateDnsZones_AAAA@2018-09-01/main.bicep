param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource AAAA 'Microsoft.Network/privateDnsZones/AAAA@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    aaaaRecords: [
      {
        ipv6Address: 'fd5d:70bc:930e:d008:0000:0000:0000:7334'
      }
      {
        ipv6Address: 'fd5d:70bc:930e:d008::7335'
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

