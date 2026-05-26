param location string = 'westus'
param resource_name string = 'acctest0001'

resource spacecraft 'Microsoft.Orbital/spacecrafts@2022-11-01' = {
  location: location
  name: resource_name
  properties: {
    links: [
      {
        bandwidthMHz: 100
        centerFrequencyMHz: 101
        direction: 'Uplink'
        name: 'linkname'
        polarization: 'LHCP'
      }
    ]
    noradId: '12345'
    titleLine: 'AQUA'
    tleLine1: '1 23455U 94089A   97320.90946019  .00000140  00000-0  10191-3 0  2621'
    tleLine2: '2 23455  99.0090 272.6745 0008546 223.1686 136.8816 14.11711747148495'
  }
}

