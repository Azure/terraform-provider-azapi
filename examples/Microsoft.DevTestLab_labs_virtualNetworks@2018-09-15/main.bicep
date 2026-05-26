param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource lab 'Microsoft.DevTestLab/labs@2018-09-15' = {
  location: location
  name: resource_name
  properties: {
    labStorageType: 'Premium'
  }
}

resource virtualNetwork 'Microsoft.DevTestLab/labs/virtualNetworks@2018-09-15' = {
  parent: lab
  name: resource_name
  properties: {
    description: ''
    subnetOverrides: [
      {
        labSubnetName: data.azapi_resource_id.subnet.name
        resourceId: data.azapi_resource_id.subnet.id
        useInVmCreationPermission: 'Allow'
        usePublicIpAddressPermission: 'Allow'
      }
    ]
  }
}

