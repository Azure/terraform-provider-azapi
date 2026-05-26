param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource server 'Microsoft.AnalysisServices/servers@2017-08-01' = {
  location: location
  name: resource_name
  properties: {
    asAdministrators: {
      members: []
    }
    ipV4FirewallSettings: {
      enablePowerBIService: false
      firewallRules: []
    }
  }
  sku: {
    name: 'B1'
  }
}

