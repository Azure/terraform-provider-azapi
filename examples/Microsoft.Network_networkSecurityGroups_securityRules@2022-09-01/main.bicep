param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource networkSecurityGroup 'Microsoft.Network/networkSecurityGroups@2022-07-01' = {
  location: location
  name: 'mi-security-group1-230630034008554952'
  properties: {
    securityRules: []
  }
}

resource securityRule 'Microsoft.Network/networkSecurityGroups/securityRules@2022-09-01' = {
  parent: networkSecurityGroup
  name: 'allow_management_inbound'
  properties: {
    access: 'Allow'
    destinationAddressPrefix: '*'
    destinationPortRange: ''
    destinationPortRanges: [
      '9000'
      '1438'
      '1440'
      '9003'
      '1452'
    ]
    direction: 'Inbound'
    priority: 106
    protocol: 'Tcp'
    sourceAddressPrefix: '*'
    sourcePortRange: '*'
  }
}

