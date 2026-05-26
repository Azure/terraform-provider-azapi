param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource N 'Microsoft.Network/dnsZones/NS@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    NSRecords: [
      {
        nsdname: 'ns1.contoso.com'
      }
      {
        nsdname: 'ns2.contoso.com'
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

