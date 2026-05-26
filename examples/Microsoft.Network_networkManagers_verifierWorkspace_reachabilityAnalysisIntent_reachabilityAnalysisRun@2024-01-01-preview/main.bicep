param location string = 'westeurope'
param resource_name string = 'acctest0001'
param vm_admin_password string = null

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

resource networkManager 'Microsoft.Network/networkManagers@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    description: ''
    networkManagerScopeAccesses: [
      'SecurityAdmin'
    ]
    networkManagerScopes: {
      managementGroups: []
      subscriptions: [
        '/subscriptions/${data.azapi_client_config.current.subscription_id}'
      ]
    }
  }
}

resource reachabilityAnalysisIntent 'Microsoft.Network/networkManagers/verifierWorkspaces/reachabilityAnalysisIntents@2024-01-01-preview' = {
  parent: verifierWorkspace
  name: resource_name
  properties: {
    description: 'A sample reachability analysis intent'
    destinationResourceId: virtualMachine.id
    ipTraffic: {
      destinationIps: [
        '10.4.0.1'
      ]
      destinationPorts: [
        '0'
      ]
      protocols: [
        'Any'
      ]
      sourceIps: [
        '10.4.0.0'
      ]
      sourcePorts: [
        '0'
      ]
    }
    sourceResourceId: virtualMachine.id
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

resource verifierWorkspace 'Microsoft.Network/networkManagers/verifierWorkspaces@2024-01-01-preview' = {
  parent: networkManager
  location: location
  name: resource_name
  tags: {
    myTag: 'testTag'
  }
  properties: {
    description: 'A sample workspace'
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

