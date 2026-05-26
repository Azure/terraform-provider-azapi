param location string = 'eastus'
param resource_name string = 'acctest0001'

resource networkInterface 'Microsoft.Network/networkInterfaces@2024-05-01' = {
  location: location
  name: '${resource_name}-nic'
  properties: {
    enableAcceleratedNetworking: false
    enableIPForwarding: false
    ipConfigurations: [
      {
        name: 'internal'
        properties: {
          primary: false
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

resource runCommand 'Microsoft.Compute/virtualMachines/runCommands@2023-03-01' = {
  parent: virtualMachine
  location: location
  name: '${resource_name}-runcommand'
  properties: {
    asyncExecution: false
    errorBlobUri: ''
    outputBlobUri: ''
    parameters: []
    protectedParameters: []
    runAsPassword: ''
    runAsUser: ''
    source: {
      script: 'echo \'hello world\''
    }
    timeoutInSeconds: 1200
    treatFailureAsDeploymentFailure: true
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2024-05-01' = {
  parent: virtualNetwork
  name: 'internal'
  properties: {
    addressPrefix: '10.0.2.0/24'
    defaultOutboundAccess: true
    delegations: []
    privateEndpointNetworkPolicies: 'Disabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource userAssignedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  location: location
  name: '${resource_name}-uai'
}

resource virtualMachine 'Microsoft.Compute/virtualMachines@2024-03-01' = {
  identity: [
    {
      identity_ids: [
        userAssignedIdentity.id
      ]
      type: 'SystemAssigned, UserAssigned'
    }
  ]
  location: location
  name: '${resource_name}-vm'
  properties: {
    additionalCapabilities: {}
    applicationProfile: {
      galleryApplications: []
    }
    diagnosticsProfile: {
      bootDiagnostics: {
        enabled: false
        storageUri: ''
      }
    }
    extensionsTimeBudget: 'PT1H30M'
    hardwareProfile: {
      vmSize: 'Standard_B2s'
    }
    networkProfile: {
      networkInterfaces: [
        {
          id: networkInterface.id
          properties: {
            primary: true
          }
        }
      ]
    }
    osProfile: {
      adminPassword: 'Pa-tn93e'
      adminUsername: 'adminuser'
      allowExtensionOperations: true
      computerName: '${resource_name}-vm'
      linuxConfiguration: {
        disablePasswordAuthentication: false
        patchSettings: {
          assessmentMode: 'ImageDefault'
          patchMode: 'ImageDefault'
        }
        provisionVMAgent: true
        ssh: {
          publicKeys: []
        }
      }
      secrets: []
    }
    priority: 'Regular'
    storageProfile: {
      dataDisks: []
      imageReference: {
        offer: '0001-com-ubuntu-server-jammy'
        publisher: 'Canonical'
        sku: '22_04-lts'
        version: 'latest'
      }
      osDisk: {
        caching: 'ReadWrite'
        createOption: 'FromImage'
        managedDisk: {
          storageAccountType: 'Premium_LRS'
        }
        osType: 'Linux'
        writeAcceleratorEnabled: false
      }
    }
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  location: location
  name: '${resource_name}-vnet'
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    privateEndpointVNetPolicies: 'Disabled'
    subnets: []
  }
}

