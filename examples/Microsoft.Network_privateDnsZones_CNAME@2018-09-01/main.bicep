param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource CNAME 'Microsoft.Network/privateDnsZones/CNAME@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    cnameRecord: {
      cname: 'contoso.com'
    }
    metadata: {}
    ttl: 300
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

