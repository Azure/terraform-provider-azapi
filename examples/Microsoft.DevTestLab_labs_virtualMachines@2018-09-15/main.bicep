param location string = 'westeurope'
param resource_name string = 'acctest0001'
param vm_password string = null

resource lab 'Microsoft.DevTestLab/labs@2018-09-15' = {
  location: location
  name: resource_name
  properties: {
    labStorageType: 'Premium'
  }
}

resource virtualMachine 'Microsoft.DevTestLab/labs/virtualMachines@2018-09-15' = {
  parent: lab
  location: location
  name: resource_name
  properties: {
    allowClaim: true
    disallowPublicIpAddress: false
    galleryImageReference: {
      offer: 'WindowsServer'
      osType: 'Windows'
      publisher: 'MicrosoftWindowsServer'
      sku: '2012-Datacenter'
      version: 'latest'
    }
    isAuthenticationWithSshKey: false
    labSubnetName: data.azapi_resource_id.subnet.name
    labVirtualNetworkId: virtualNetwork.id
    networkInterface: {}
    notes: ''
    osType: 'Windows'
    password: vm_password
    size: 'Standard_F2'
    storageType: 'Standard'
    userName: 'acct5stU5er'
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

