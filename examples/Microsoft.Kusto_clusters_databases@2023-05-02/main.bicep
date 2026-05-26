param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.Kusto/clusters@2023-05-02' = {
  location: location
  name: resource_name
  properties: {
    enableAutoStop: true
    enableDiskEncryption: false
    enableDoubleEncryption: false
    enablePurge: false
    enableStreamingIngest: false
    engineType: 'V2'
    publicIPType: 'IPv4'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    trustedExternalTenants: []
  }
  sku: {
    capacity: 1
    name: 'Dev(No SLA)_Standard_D11_v2'
    tier: 'Basic'
  }
}

resource database 'Microsoft.Kusto/clusters/databases@2023-05-02' = {
  parent: cluster
  location: location
  name: resource_name
  kind: 'ReadWrite'
  properties: {}
}

