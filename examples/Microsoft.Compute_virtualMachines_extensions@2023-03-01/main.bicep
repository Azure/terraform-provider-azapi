param location string = 'westeurope'
param resource_name string = 'acctest0001'
param vm_admin_password string = null

resource extension 'Microsoft.Compute/virtualMachines/extensions@2023-03-01' = {
  parent: virtualMachine
  location: location
  name: resource_name
  properties: {
    autoUpgradeMinorVersion: false
    enableAutomaticUpgrade: false
    publisher: 'Microsoft.Azure.Extensions'
    settings: {
      commandToExecute: 'hostname'
    }
    suppressFailures: false
    type: 'CustomScript'
    typeHandlerVersion: '2.0'
  }
  tags: {
    environment: 'Production'
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
      adminPassword: vm_admin_password
      adminUsername: 'testadmin'
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
        name: 'myosdisk1'
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

