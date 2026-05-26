param admin_password string = null
param admin_username string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource autoScaleSetting 'Microsoft.Insights/autoScaleSettings@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    enabled: true
    notifications: []
    profiles: [
      {
        capacity: {
          default: '1'
          maximum: '10'
          minimum: '1'
        }
        name: 'metricRules'
        rules: [
          {
            metricTrigger: {
              dimensions: []
              dividePerInstance: true
              metricName: 'Percentage CPU'
              metricNamespace: ''
              metricResourceUri: virtualMachineScaleSet.id
              operator: 'GreaterThan'
              statistic: 'Average'
              threshold: 75
              timeAggregation: 'Last'
              timeGrain: 'PT1M'
              timeWindow: 'PT5M'
            }
            scaleAction: {
              cooldown: 'PT1M'
              direction: 'Increase'
              type: 'ChangeCount'
              value: '1'
            }
          }
        ]
      }
    ]
    targetResourceUri: virtualMachineScaleSet.id
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: 'internal'
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
  name: resource_name
  properties: {
    additionalCapabilities: {}
    doNotRunExtensionsOnOverprovisionedVMs: false
    orchestrationMode: 'Uniform'
    overprovision: true
    scaleInPolicy: {
      forceDeletion: false
      rules: [
        'Default'
      ]
    }
    singlePlacementGroup: true
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
        networkInterfaceConfigurations: [
          {
            name: 'TestNetworkProfile-230630033559396108'
            properties: {
              dnsSettings: {
                dnsServers: []
              }
              enableAcceleratedNetworking: false
              enableIPForwarding: false
              ipConfigurations: [
                {
                  name: 'TestIPConfiguration'
                  properties: {
                    applicationGatewayBackendAddressPools: []
                    applicationSecurityGroups: []
                    loadBalancerBackendAddressPools: []
                    loadBalancerInboundNatPools: []
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
        adminPassword: admin_password
        adminUsername: admin_username
        computerNamePrefix: 'testvm-230630033559396108'
        linuxConfiguration: {
          disablePasswordAuthentication: false
          provisionVMAgent: true
          ssh: {
            publicKeys: [
              {
                keyData: 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com'
                path: '/home/myadmin/.ssh/authorized_keys'
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
            storageAccountType: 'StandardSSD_LRS'
          }
          osType: 'Linux'
          writeAcceleratorEnabled: false
        }
      }
    }
  }
  sku: {
    capacity: 2
    name: 'Standard_F2'
    tier: 'Standard'
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

