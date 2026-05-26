param location string = 'centralus'
param resource_name string = 'acctest0001'

resource authorization 'Microsoft.AVS/privateClouds/authorizations@2022-05-01' = {
  parent: privateCloud
  name: resource_name
}

resource privateCloud 'Microsoft.AVS/privateClouds@2022-05-01' = {
  location: location
  name: resource_name
  properties: {
    internet: 'Disabled'
    managementCluster: {
      clusterSize: 3
    }
    networkBlock: '192.168.48.0/22'
  }
  sku: {
    name: 'av36'
  }
}

