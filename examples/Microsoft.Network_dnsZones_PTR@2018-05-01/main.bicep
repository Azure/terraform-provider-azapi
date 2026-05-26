param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource PTR 'Microsoft.Network/dnsZones/PTR@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    PTRRecords: [
      {
        ptrdname: 'hashicorp.com'
      }
      {
        ptrdname: 'microsoft.com'
      }
    ]
    TTL: 300
    metadata: {}
  }
}

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

