param admin_password string = null
param admin_username string = null
param attached_resource_name string = 'acctest0002'
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource attachedManagedDisk 'Microsoft.Compute/disks@2023-10-02' = {
  location: location
  name: local.attached_os_disk_name
  properties: {
    creationData: {
      createOption: 'Copy'
      sourceResourceId: snapshot.id
    }
    diskSizeGB: 30
    encryption: {
      type: 'EncryptionAtRestWithPlatformKey'
    }
    hyperVGeneration: 'V1'
    networkAccessPolicy: 'AllowAll'
    osType: 'Linux'
    publicNetworkAccess: 'Enabled'
    supportedCapabilities: {
      architecture: 'x64'
    }
  }
  sku: {
    name: 'Standard_LRS'
  }
  zones: [
    '1'
  ]
}

resource attachedNetworkInterface 'Microsoft.Network/networkInterfaces@2022-07-01' = {
  location: location
  name: attached_resource_name
  properties: {
    enableAcceleratedNetworking: false
    enableIPForwarding: false
    ipConfigurations: [
      {
        name: 'testconfiguration2'
        properties: {
          primary: true
          privateIPAddressVersion: 'IPv4'
          privateIPAllocationMethod: 'Dynamic'
          subnet: {
            id: subnet.id
          }
        }
      }
    ]
  }
}

resource attachedVirtualMachine 'Microsoft.Compute/virtualMachines@2023-03-01' = {
  location: location
  name: attached_resource_name
  properties: {
    hardwareProfile: {
      vmSize: 'Standard_F2'
    }
    networkProfile: {
      networkInterfaces: [
        {
          id: attachedNetworkInterface.id
          properties: {
            primary: false
          }
        }
      ]
    }
    storageProfile: {
      osDisk: {
        caching: 'ReadWrite'
        createOption: 'Attach'
        managedDisk: {
          id: attachedManagedDisk.id
        }
        name: local.attached_os_disk_name
        osType: 'Linux'
        writeAcceleratorEnabled: false
      }
    }
  }
}

resource networkInterface 'Microsoft.Network/networkInterfaces@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    enableAcceleratedNetworking: false
    enableIPForwarding: false
    ipConfigurations: [
      {
        name: 'testconfiguration1'
        properties: {
          primary: true
          privateIPAddressVersion: 'IPv4'
          privateIPAllocationMethod: 'Dynamic'
          subnet: {
            id: subnet.id
          }
        }
      }
    ]
  }
}

resource snapshot 'Microsoft.Compute/snapshots@2023-10-02' = {
  location: location
  name: resource_name
  properties: {
    creationData: {
      createOption: 'Copy'
      sourceResourceId: data.azapi_resource.managedDisk.id
    }
    diskSizeGB: 30
    encryption: {
      type: 'EncryptionAtRestWithPlatformKey'
    }
    hyperVGeneration: 'V1'
    incremental: true
    networkAccessPolicy: 'AllowAll'
    osType: 'Linux'
    publicNetworkAccess: 'Enabled'
    supportedCapabilities: {
      architecture: 'x64'
    }
  }
  sku: {
    name: 'Standard_ZRS'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.0.2.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource virtualMachine 'Microsoft.Compute/virtualMachines@2023-03-01' = {
  location: location
  name: resource_name
  properties: {
    hardwareProfile: {
      vmSize: 'Standard_F2'
    }
    networkProfile: {
      networkInterfaces: [
        {
          id: networkInterface.id
          properties: {
            primary: false
          }
        }
      ]
    }
    osProfile: {
      adminPassword: admin_password
      adminUsername: admin_username
      computerName: 'hostname230630032848831819'
      linuxConfiguration: {
        disablePasswordAuthentication: false
      }
    }
    storageProfile: {
      imageReference: {
        offer: 'UbuntuServer'
        publisher: 'Canonical'
        sku: '16.04-LTS'
        version: 'latest'
      }
      osDisk: {
        caching: 'ReadWrite'
        createOption: 'FromImage'
        name: local.os_disk_name
        writeAcceleratorEnabled: false
      }
    }
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

