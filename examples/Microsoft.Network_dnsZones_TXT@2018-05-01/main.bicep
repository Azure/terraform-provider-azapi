param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource TXT 'Microsoft.Network/dnsZones/TXT@2018-05-01' = {
  parent: dnsZone
  name: resource_name
  properties: {
    TTL: 300
    TXTRecords: [
      {
        value: [
          'Quick brown fox'
        ]
      }
      {
        value: [
          'A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text.....'
          '.A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text....'
          '..A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......'
        ]
      }
    ]
    metadata: {}
  }
}

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

