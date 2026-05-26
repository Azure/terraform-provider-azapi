param location string = 'eastus'
param resource_name string = 'acctest0001'

resource dataBoxEdgeDevice 'Microsoft.DataBoxEdge/dataBoxEdgeDevices@2022-03-01' = {
  location: location
  name: resource_name
  sku: {
    name: 'EdgeP_Base'
    tier: 'Standard'
  }
}

