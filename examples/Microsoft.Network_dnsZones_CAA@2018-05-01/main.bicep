param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource CAA 'Microsoft.Network/dnsZones/CAA@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    TTL: 300
    caaRecords: [
      {
        flags: 1
        tag: 'issuewild'
        value: ';'
      }
      {
        flags: 0
        tag: 'iodef'
        value: 'mailto:terraform@nonexist.tld'
      }
      {
        flags: 0
        tag: 'issue'
        value: 'example.com'
      }
      {
        flags: 0
        tag: 'issue'
        value: 'example.net'
      }
    ]
    metadata: {}
  }
}

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

