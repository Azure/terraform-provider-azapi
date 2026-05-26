param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.ServiceFabric/clusters@2021-06-01' = {
  location: location
  name: resource_name
  properties: {
    addOnFeatures: []
    fabricSettings: []
    managementEndpoint: 'http://example:80'
    nodeTypes: [
      {
        capacities: {}
        clientConnectionEndpointPort: 2020
        durabilityLevel: 'Bronze'
        httpGatewayEndpointPort: 80
        isPrimary: true
        isStateless: false
        multipleAvailabilityZones: false
        name: 'first'
        placementProperties: {}
        vmInstanceCount: 3
      }
    ]
    reliabilityLevel: 'Bronze'
    upgradeMode: 'Automatic'
    vmImage: 'Windows'
  }
}

