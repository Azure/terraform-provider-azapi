param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource SRV 'Microsoft.Network/privateDnsZones/SRV@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    metadata: {}
    srvRecords: [
      {
        port: 8080
        priority: 10
        target: 'target2.contoso.com'
        weight: 10
      }
      {
        port: 8080
        priority: 1
        target: 'target1.contoso.com'
        weight: 5
      }
    ]
    ttl: 300
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

