param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource MX 'Microsoft.Network/privateDnsZones/MX@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    metadata: {}
    mxRecords: [
      {
        exchange: 'mx1.contoso.com'
        preference: 10
      }
      {
        exchange: 'mx2.contoso.com'
        preference: 10
      }
    ]
    ttl: 300
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

