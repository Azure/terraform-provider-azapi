param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource TXT 'Microsoft.Network/privateDnsZones/TXT@2018-09-01' = {
  parent: privateDnsZone
  name: resource_name
  properties: {
    metadata: {}
    ttl: 300
    txtRecords: [
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
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2018-09-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

