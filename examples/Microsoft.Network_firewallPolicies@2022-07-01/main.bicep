param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource firewallPolicy 'Microsoft.Network/firewallPolicies@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    threatIntelMode: 'Alert'
  }
}

