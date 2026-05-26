param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource agentPool 'Microsoft.ContainerService/managedClusters/agentPools@2023-04-02-preview' = {
  parent: managedCluster
  name: 'internal'
  properties: {
    count: 1
    mode: 'User'
    vmSize: 'Standard_DS2_v2'
  }
}

resource managedCluster 'Microsoft.ContainerService/managedClusters@2023-04-02-preview' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    agentPoolProfiles: [
      {
        count: 1
        mode: 'System'
        name: 'default'
        vmSize: 'Standard_DS2_v2'
      }
    ]
    dnsPrefix: resource_name
  }
}

