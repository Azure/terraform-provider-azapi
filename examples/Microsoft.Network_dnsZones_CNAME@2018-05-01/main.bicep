param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource CNAME 'Microsoft.Network/dnsZones/CNAME@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    CNAMERecord: {
      cname: '${resource_name}.webpubsub.azure.com'
    }
    TTL: 3600
    metadata: {}
    targetResource: {}
  }
}

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

