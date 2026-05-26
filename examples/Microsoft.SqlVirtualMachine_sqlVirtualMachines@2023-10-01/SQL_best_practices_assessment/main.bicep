param admin_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataCollectionEndpoint 'Microsoft.Insights/dataCollectionEndpoints@2022-06-01' = {
  location: resourceGroup().location
  name: '${location}-DCE-1'
  properties: {
    networkAcls: {
      publicNetworkAccess: 'Enabled'
    }
  }
}

resource dataCollectionRule 'Microsoft.Insights/dataCollectionRules@2022-06-01' = {
  location: resourceGroup().location
  name: '${workspace.properties.customerId}_${resourceGroup().location}_DCR_1'
  properties: {
    dataCollectionEndpointId: dataCollectionEndpoint.id
    dataFlows: [
      {
        destinations: [
          workspace.name
        ]
        outputStream: 'Custom-SqlAssessment_CL'
        streams: [
          'Custom-SqlAssessment_CL'
        ]
        transformKql: 'source'
      }
    ]
    dataSources: {
      logFiles: [
        {
          filePatterns: [
            'C:\\Windows\\System32\\config\\systemprofile\\AppData\\Local\\Microsoft SQL Server IaaS Agent\\Assessment\\*.csv'
          ]
          format: 'text'
          name: 'Custom-SqlAssessment_CL'
          settings: {
            text: {
              recordStartTimestampFormat: 'ISO 8601'
            }
          }
          streams: [
            'Custom-SqlAssessment_CL'
          ]
        }
      ]
    }
    description: ''
    destinations: {
      logAnalytics: [
        {
          name: workspace.name
          workspaceResourceId: workspace.id
        }
      ]
    }
    streamDeclarations: {
      'Custom-SqlAssessment_CL': {
        columns: [
          {
            name: 'TimeGenerated'
            type: 'datetime'
          }
          {
            name: 'RawData'
            type: 'string'
          }
        ]
      }
    }
  }
}

resource dataCollectionRuleAssociation 'Microsoft.Insights/dataCollectionRuleAssociations@2022-06-01' = {
  parent: virtualMachine
  name: '${workspace.properties.customerId}_${resourceGroup().location}_DCRA_1'
  properties: {
    dataCollectionRuleId: dataCollectionRule.id
  }
}

resource extension 'Microsoft.Compute/virtualMachines/extensions@2024-07-01' = {
  parent: virtualMachine
  location: 'westeurope'
  name: 'AzureMonitorWindowsAgent'
  properties: {
    autoUpgradeMinorVersion: true
    enableAutomaticUpgrade: true
    publisher: 'Microsoft.Azure.Monitor'
    suppressFailures: false
    type: 'AzureMonitorWindowsAgent'
    typeHandlerVersion: '1.0'
  }
}

resource networkInterface 'Microsoft.Network/networkInterfaces@2024-05-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    auxiliaryMode: 'None'
    auxiliarySku: 'None'
    disableTcpStateTracking: false
    dnsSettings: {
      dnsServers: []
    }
    enableAcceleratedNetworking: false
    enableIPForwarding: false
    ipConfigurations: [
      {
        name: 'testconfiguration1'
        properties: {
          primary: true
          privateIPAddress: '10.0.0.4'
          privateIPAddressVersion: 'IPv4'
          privateIPAllocationMethod: 'Dynamic'
          publicIPAddress: {
            id: publicIPAddress.id
          }
          subnet: {
            id: subnet.id
          }
        }
        type: 'Microsoft.Network/networkInterfaces/ipConfigurations'
      }
    ]
    nicType: 'Standard'
  }
}

resource networkSecurityGroup 'Microsoft.Network/networkSecurityGroups@2024-05-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    securityRules: [
      {
        name: 'MSSQLRule'
        properties: {
          access: 'Allow'
          destinationAddressPrefix: '*'
          destinationAddressPrefixes: []
          destinationPortRange: '1433'
          destinationPortRanges: []
          direction: 'Inbound'
          priority: 1001
          protocol: 'Tcp'
          sourceAddressPrefix: '167.220.255.0/25'
          sourceAddressPrefixes: []
          sourcePortRange: '*'
          sourcePortRanges: []
        }
      }
    ]
  }
}

resource publicIPAddress 'Microsoft.Network/publicIPAddresses@2024-05-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    ddosSettings: {
      protectionMode: 'VirtualNetworkInherited'
    }
    idleTimeoutInMinutes: 4
    ipTags: []
    publicIPAddressVersion: 'IPv4'
    publicIPAllocationMethod: 'Dynamic'
  }
  sku: {
    name: 'Basic'
    tier: 'Regional'
  }
}

resource sqlvirtualMachine 'Microsoft.SqlVirtualMachine/sqlVirtualMachines@2023-10-01' = {
  location: virtualMachine.location
  name: virtualMachine.name
  properties: {
    assessmentSettings: {
      enable: true
      runImmediately: false
      schedule: {
        dayOfWeek: 'Monday'
        enable: true
        startTime: '00:00'
        weeklyInterval: 1
      }
    }
    enableAutomaticUpgrade: true
    leastPrivilegeMode: 'Enabled'
    sqlImageOffer: 'SQL2017-WS2016'
    sqlImageSku: 'Developer'
    sqlManagement: 'Full'
    sqlServerLicenseType: 'PAYG'
    virtualMachineResourceId: virtualMachine.id
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.0.0.0/24'
    networkSecurityGroup: {
      id: networkSecurityGroup.id
    }
  }
}

resource table 'Microsoft.OperationalInsights/workspaces/tables@2023-09-01' = {
  parent: workspace
  name: 'SqlAssessment_CL'
  properties: {
    schema: {
      columns: [
        {
          name: 'TimeGenerated'
          type: 'datetime'
        }
        {
          name: 'RawData'
          type: 'string'
        }
      ]
      name: 'SqlAssessment_CL'
    }
  }
}

resource virtualMachine 'Microsoft.Compute/virtualMachines@2024-07-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    hardwareProfile: {
      vmSize: 'Standard_F2s'
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
      adminUsername: 'testadmin'
      allowExtensionOperations: true
      computerName: 'winhost01'
      secrets: []
      windowsConfiguration: {
        enableAutomaticUpdates: true
        patchSettings: {
          assessmentMode: 'ImageDefault'
          patchMode: 'AutomaticByOS'
        }
        provisionVMAgent: true
        timeZone: 'Pacific Standard Time'
      }
    }
    storageProfile: {
      dataDisks: []
      imageReference: {
        offer: 'SQL2017-WS2016'
        publisher: 'MicrosoftSQLServer'
        sku: 'SQLDEV'
        version: 'latest'
      }
      osDisk: {
        caching: 'ReadOnly'
        createOption: 'FromImage'
        deleteOption: 'Detach'
        diskSizeGB: 127
        managedDisk: {
          storageAccountType: 'Premium_LRS'
        }
        name: 'acctvm-250116171212663925OSDisk'
        osType: 'Windows'
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
  }
}

resource workspace 'Microsoft.OperationalInsights/workspaces@2020-08-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    features: {
      disableLocalAuth: false
      enableLogAccessUsingOnlyResourcePermissions: true
      legacy: 0
      searchVersion: 1
    }
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
    retentionInDays: 30
    sku: {
      name: 'PerGB2018'
    }
    workspaceCapping: {
      dailyQuotaGb: -1
    }
  }
}

