param location string = 'eastus'
param resource_name string = 'acctest0001'

resource standbyVirtualMachinePool 'Microsoft.StandbyPool/standbyVirtualMachinePools@2025-03-01' = {
  location: location
  name: '${resource_name}-pool'
  properties: {
    attachedVirtualMachineScaleSetId: virtualMachineScaleSet.id
    elasticityProfile: {
      maxReadyCapacity: 2
      minReadyCapacity: 1
    }
    virtualMachineState: 'Running'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: '${resource_name}-subnet'
  properties: {
    addressPrefix: '10.0.2.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource virtualMachineScaleSet 'Microsoft.Compute/virtualMachineScaleSets@2023-03-01' = {
  location: location
  name: '${resource_name}-vmss'
  properties: {
    additionalCapabilities: {}
    orchestrationMode: 'Flexible'
    platformFaultDomainCount: 1
    singlePlacementGroup: false
    upgradePolicy: {
      mode: 'Manual'
    }
    virtualMachineProfile: {
      diagnosticsProfile: {
        bootDiagnostics: {
          enabled: false
          storageUri: ''
        }
      }
      extensionProfile: {
        extensionsTimeBudget: 'PT1H30M'
      }
      networkProfile: {
        networkApiVersion: '2022-07-01'
        networkInterfaceConfigurations: [
          {
            name: 'example'
            properties: {
              dnsSettings: {
                dnsServers: []
              }
              enableAcceleratedNetworking: false
              enableIPForwarding: false
              ipConfigurations: [
                {
                  name: 'internal'
                  properties: {
                    applicationGatewayBackendAddressPools: []
                    applicationSecurityGroups: []
                    loadBalancerBackendAddressPools: []
                    primary: true
                    privateIPAddressVersion: 'IPv4'
                    subnet: {
                      id: subnet.id
                    }
                  }
                }
              ]
              primary: true
            }
          }
        ]
      }
      osProfile: {
        adminUsername: 'adminuser'
        computerNamePrefix: resource_name
        linuxConfiguration: {
          disablePasswordAuthentication: true
          provisionVMAgent: true
          ssh: {
            publicKeys: [
              {
                keyData: 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com'
                path: '/home/adminuser/.ssh/authorized_keys'
              }
            ]
          }
        }
        secrets: []
      }
      priority: 'Regular'
      storageProfile: {
        dataDisks: []
        imageReference: {
          offer: 'UbuntuServer'
          publisher: 'Canonical'
          sku: '16.04-LTS'
          version: 'latest'
        }
        osDisk: {
          caching: 'ReadWrite'
          createOption: 'FromImage'
          managedDisk: {
            storageAccountType: 'Standard_LRS'
          }
          osType: 'Linux'
          writeAcceleratorEnabled: false
        }
      }
    }
  }
  sku: {
    capacity: 1
    name: 'Standard_D1_v2'
    tier: 'Standard'
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
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
    subnets: []
  }
}

