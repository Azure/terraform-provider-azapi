param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource MX 'Microsoft.Network/dnsZones/MX@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    MXRecords: [
      {
        exchange: 'mail2.contoso.com'
        preference: 20
      }
      {
        exchange: 'mail1.contoso.com'
        preference: 10
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

