param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource SRV 'Microsoft.Network/dnsZones/SRV@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    SRVRecords: [
      {
        port: 8080
        priority: 2
        target: 'target2.contoso.com'
        weight: 25
      }
      {
        port: 8080
        priority: 1
        target: 'target1.contoso.com'
        weight: 5
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

