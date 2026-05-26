param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource disk 'Microsoft.Compute/disks@2022-03-02' = {
  location: location
  name: resource_name
  properties: {
    creationData: {
      createOption: 'Empty'
    }
    diskSizeGB: 10
    encryption: {
      type: 'EncryptionAtRestWithPlatformKey'
    }
    networkAccessPolicy: 'AllowAll'
    osType: ''
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    name: 'Standard_LRS'
  }
}

