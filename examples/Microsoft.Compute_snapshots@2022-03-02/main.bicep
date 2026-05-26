param location string = 'westus'
param resource_name string = 'acctest0001'

resource disk 'Microsoft.Compute/disks@2023-04-02' = {
  location: location
  name: '${resource_name}disk'
  properties: {
    creationData: {
      createOption: 'Empty'
      performancePlus: false
    }
    diskSizeGB: 10
    encryption: {
      type: 'EncryptionAtRestWithPlatformKey'
    }
    networkAccessPolicy: 'AllowAll'
    optimizedForFrequentAttach: false
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    name: 'Standard_LRS'
  }
}

resource snapshot 'Microsoft.Compute/snapshots@2022-03-02' = {
  location: location
  name: '${resource_name}snapshot'
  properties: {
    creationData: {
      createOption: 'Copy'
      sourceUri: disk.id
    }
    diskSizeGB: 20
    incremental: false
    networkAccessPolicy: 'AllowAll'
    publicNetworkAccess: 'Enabled'
  }
}

