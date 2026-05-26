param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource PTR 'Microsoft.Network/privateDnsZones/PTR@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    metadata: {}
    ptrRecords: [
      {
        ptrdname: 'test2.contoso.com'
      }
      {
        ptrdname: 'test.contoso.com'
      }
    ]
    ttl: 300
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '230630033756174960.0.10.in-addr.arpa'
}

